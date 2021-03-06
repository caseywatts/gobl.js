package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/imdario/mergo"
	"github.com/labstack/echo/v4"
	"gopkg.in/yaml.v3"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl.cli/internal/iotools"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/schema"
)

// BuildOptions are the options to pass to the Build function.
type BuildOptions struct {
	Template   io.Reader
	Data       io.Reader
	DocType    string
	SetYAML    map[string]string
	SetString  map[string]string
	SetFile    map[string]string
	PrivateKey *dsig.PrivateKey
}

// decodeInto unmarshals in as YAML, then merges it into dest.
func decodeInto(ctx context.Context, dest *map[string]interface{}, in io.Reader) error {
	var intermediate map[string]interface{}
	dec := yaml.NewDecoder(iotools.CancelableReader(ctx, in))
	if err := dec.Decode(&intermediate); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := mergo.Merge(dest, intermediate, mergo.WithOverride); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	return nil
}

// Build builds and validates a GOBL envelope from the opts.
func Build(ctx context.Context, opts *BuildOptions) (*gobl.Envelope, error) {
	encoded, err := prepareIntermediate(ctx, opts, docInEnvelopeSchemaData)
	if err != nil {
		return nil, err
	}

	// Prepare an empty envelope as we assume the consumer is providing one already.
	env := new(gobl.Envelope)
	if err := json.Unmarshal(encoded, env); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := finalizeEnvelope(ctx, env, opts); err != nil {
		return nil, err
	}
	return env, nil
}

// Envelop assumes the incoming BuildOptions define the contents of a document
// payload and we need to prepare the envelope around it.
func Envelop(ctx context.Context, opts *BuildOptions) (*gobl.Envelope, error) {
	encoded, err := prepareIntermediate(ctx, opts, docSchemaData)
	if err != nil {
		return nil, err
	}

	// Prepare a new envelope as the intention is to insert the encoded data
	env := gobl.NewEnvelope()
	if err = json.Unmarshal(encoded, env.Document); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := finalizeEnvelope(ctx, env, opts); err != nil {
		return nil, err
	}
	return env, nil
}

type schemaDataCB func(schema.ID) map[string]interface{}

func prepareIntermediate(ctx context.Context, opts *BuildOptions, schemaDataFunc schemaDataCB) ([]byte, error) {
	values, err := parseSets(opts)
	if err != nil {
		return nil, err
	}
	var intermediate map[string]interface{}

	if opts.Template != nil {
		if err = decodeInto(ctx, &intermediate, opts.Template); err != nil {
			return nil, err
		}
	}
	if err = decodeInto(ctx, &intermediate, opts.Data); err != nil {
		return nil, err
	}

	if err = mergo.Merge(&intermediate, values, mergo.WithOverride); err != nil {
		return nil, echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	if opts.DocType != "" {
		schema := FindType(opts.DocType)
		if schema == "" {
			return nil, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("unrecognized doc type: %q", opts.DocType))
		}
		if err = mergo.Merge(&intermediate, schemaDataFunc(schema)); err != nil {
			return nil, err
		}
	}

	return json.Marshal(intermediate)
}

func finalizeEnvelope(ctx context.Context, env *gobl.Envelope, opts *BuildOptions) error {
	if len(env.Signatures) > 0 {
		return echo.NewHTTPError(http.StatusConflict, "document has already been signed")
	}
	if err := env.Complete(); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	if opts.PrivateKey == nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "signing key required")
	}
	if !env.Head.Draft {
		if err := env.Sign(opts.PrivateKey); err != nil {
			return err
		}
	}
	return nil
}

func docInEnvelopeSchemaData(schema schema.ID) map[string]interface{} {
	return map[string]interface{}{
		"doc": docSchemaData(schema),
	}
}

func docSchemaData(schema schema.ID) map[string]interface{} {
	return map[string]interface{}{
		"$schema": schema,
	}

}

func parseSets(opts *BuildOptions) (map[string]interface{}, error) {
	values := map[string]interface{}{}
	keys := make([]string, 0, len(opts.SetYAML))
	for k := range opts.SetYAML {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := opts.SetYAML[k]
		var parsed interface{}
		if err := yaml.Unmarshal([]byte(v), &parsed); err != nil {
			return nil, err
		}
		if err := setValue(&values, k, parsed); err != nil {
			return nil, err
		}
	}

	keys = make([]string, 0, len(opts.SetString))
	for k := range opts.SetString {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := opts.SetString[k]
		if err := setValue(&values, k, v); err != nil {
			return nil, err
		}
	}

	keys = make([]string, 0, len(opts.SetFile))
	for k := range opts.SetFile {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := opts.SetFile[k]
		f, err := os.Open(v)
		if err != nil {
			return nil, err
		}
		defer f.Close() // nolint:errcheck
		dec := yaml.NewDecoder(f)
		var val interface{}
		if err := dec.Decode(&val); err != nil {
			return nil, err
		}
		if err := setValue(&values, k, val); err != nil {
			return nil, err
		}
	}
	return values, nil
}

func setValue(values *map[string]interface{}, key string, value interface{}) error {
	key = strings.ReplaceAll(key, `\.`, "\x00")

	// If the key starts with '.', we treat that as the root of the
	// target object
	if key == "." {
		return mergo.Merge(values, value, mergo.WithOverride)
	}
	if len(key) > 1 && key[0] == '.' {
		key = key[1:]
	}

	for {
		i := strings.LastIndex(key, ".")
		if i == -1 {
			break
		}
		value = map[string]interface{}{
			strings.ReplaceAll(key[i+1:], "\x00", "."): value,
		}
		key = key[:i]
	}
	return mergo.Merge(values, map[string]interface{}{
		strings.ReplaceAll(key, "\x00", "."): value,
	}, mergo.WithOverride)
}
