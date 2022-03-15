package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"gitlab.com/flimzy/testy"
)

func Test_build_args(t *testing.T) {
	tests := []struct {
		name string
		args []string
		err  string
	}{
		{
			name: "no args",
		},
		{
			name: "invalid flag",
			args: []string{"--foo"},
			err:  `unknown flag: --foo`,
		},
		{
			name: "force long",
			args: []string{"--force"},
		},
		{
			name: "force short",
			args: []string{"-f"},
		},
		{
			name: "in-place long",
			args: []string{"--in-place"},
		},
		{
			name: "in-place short",
			args: []string{"-w"},
		},
		{
			name: "set values",
			args: []string{"--set", "foo=bar", "--set", "bar=baz", "--set", "foo=qux"},
		},
		{
			name: "set files",
			args: []string{"--set-file", "foo=foo.json"},
		},
		{
			name: "set string values",
			args: []string{"--set-string", "foo=foo", "--set-string", "bar=1234"},
		},
		{
			name: "template",
			args: []string{"--template", "foo.yaml"},
		},
		{
			name: "type",
			args: []string{"--type", "bill.Invoice"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			opts := build()

			cmd := opts.cmd()
			err := cmd.ParseFlags(tt.args)
			if tt.err == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, tt.err)
			}
			if err != nil {
				return
			}
			if d := testy.DiffInterface(testy.Snapshot(t), opts); d != nil {
				t.Error(d)
			}
		})
	}
}

func Test_build(t *testing.T) {
	readFile := func(t *testing.T, filename string) io.Reader {
		t.Helper()
		f, err := os.Open(filename)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			_ = f.Close()
		})
		return f
	}
	noTotals := func(t *testing.T) io.Reader {
		return readFile(t, "testdata/nototals.json")
	}

	tmpdir := testy.CopyTempDir(t, "testdata", 0)
	t.Cleanup(func() {
		_ = os.RemoveAll(tmpdir)
	})

	tests := []struct {
		name   string
		opts   *buildOpts
		in     io.Reader
		args   []string
		err    string
		target string
	}{
		{
			name: "invalid yaml value on command line",
			opts: &buildOpts{
				set:            map[string]string{"foo": ":"},
				privateKeyFile: "testdata/id_es256",
			},
			err: `yaml: did not find expected key`,
		},
		{
			name: "valid yaml on command line",
			in:   noTotals(t),
			opts: &buildOpts{
				set: map[string]string{
					"doc.supplier.name": "one two three",
				},
				privateKeyFile: "testdata/id_es256",
			},
		},
		{
			name: "valid string",
			in:   noTotals(t),
			opts: &buildOpts{
				setStrings: map[string]string{
					"doc.supplier.name": "123",
				},
				privateKeyFile: "testdata/id_es256",
			},
		},
		{
			name: "missing file",
			opts: &buildOpts{
				setFiles: map[string]string{
					"foo": "missing.yaml",
				},
				privateKeyFile: "testdata/id_es256",
			},
			err: `open missing.yaml: no such file or directory`,
		},
		{
			name: "valid file",
			in:   noTotals(t),
			opts: &buildOpts{
				setFiles: map[string]string{
					"doc.supplier": "testdata/supplier.yaml",
				},
				privateKeyFile: "testdata/id_es256",
			},
		},
		{
			name: "invalid stdin",
			in:   strings.NewReader("this isn't JSON"),
			opts: &buildOpts{
				privateKeyFile: "testdata/id_es256",
			},
			err: "code=400, message=yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `this is...` into map[string]interface {}",
		},
		{
			name: "success",
			in:   noTotals(t),
			opts: &buildOpts{
				privateKeyFile: "testdata/id_es256",
			},
		},
		{
			name: "no document",
			in: strings.NewReader(`{
				"head": {
					"uuid": "9d8eafd5-77be-11ec-b485-5405db9a3e49",
					"typ": "duck",
					"rgn": "ES",
					"dig": {
						"alg": "sha256",
						"val": "dce3bc3c8bf28f3d209f783917b3082ddc0339a66e9ba3aa63849e4357db1422"
					}
				},
			}`),
			opts: &buildOpts{
				privateKeyFile: "testdata/id_es256",
			},
			err: "code=422, message=no document included",
		},
		{
			name: "invalid doc",
			in: strings.NewReader(`{
				"head": {
					"uuid": "9d8eafd5-77be-11ec-b485-5405db9a3e49",
					"dig": {
						"alg": "sha256",
						"val": "dce3bc3c8bf28f3d209f783917b3082ddc0339a66e9ba3aa63849e4357db1422"
					}
				},
				doc: "foo bar baz"
			}`),
			opts: &buildOpts{
				privateKeyFile: "testdata/id_es256",
			},
			err: "code=400, message=json: cannot unmarshal string into Go struct field Envelope.doc of type gobl.schemaDoc",
		},
		{
			name: "incomplete",
			in: strings.NewReader(`{
				"head": {
					"uuid": "9d8eafd5-77be-11ec-b485-5405db9a3e49",
					"dig": {
						"alg": "sha256",
						"val": "dce3bc3c8bf28f3d209f783917b3082ddc0339a66e9ba3aa63849e4357db1422"
					}
				},
				doc: {}
			}`),
			opts: &buildOpts{
				privateKeyFile: "testdata/id_es256",
			},
			err: "code=400, message=marshal: unregistered schema: ",
		},
		{
			name: "input file",
			args: []string{"testdata/success.json"},
			opts: &buildOpts{
				privateKeyFile: "testdata/id_es256",
			},
		},
		{
			name: "recalculate",
			args: []string{"testdata/nototals.json"},
			opts: &buildOpts{
				privateKeyFile: "testdata/id_es256",
			},
		},
		{
			name: "output file",
			args: []string{"testdata/success.json", filepath.Join(tmpdir, "output-file.json")},
			opts: &buildOpts{
				privateKeyFile: "testdata/id_es256",
			},
			target: filepath.Join(tmpdir, "output-file.json"),
		},
		{
			name: "explicit stdout",
			args: []string{"testdata/success.json", "-"},
			opts: &buildOpts{
				privateKeyFile: "testdata/id_es256",
			},
		},
		{
			name: "output file exists",
			args: []string{"testdata/success.json", filepath.Join(tmpdir, "exists.json")},
			err:  "open " + tmpdir + "/exists.json: file exists",
		},
		{
			name: "overwrite output file",
			opts: &buildOpts{
				overwriteOutputFile: true,
				privateKeyFile:      "testdata/id_es256",
			},
			args:   []string{"testdata/success.json", filepath.Join(tmpdir, "overwrite.json")},
			target: filepath.Join(tmpdir, "overwrite.json"),
		},
		{
			name: "overwrite input file",
			opts: &buildOpts{
				inPlace:        true,
				privateKeyFile: "testdata/id_es256",
			},
			args:   []string{filepath.Join(tmpdir, "input.json")},
			target: filepath.Join(tmpdir, "input.json"),
		},
		{
			name: "overwrite stdin",
			opts: &buildOpts{
				inPlace: true,
			},
			err: "cannot overwrite STDIN",
		},
		{
			name: "merge values",
			opts: &buildOpts{
				set:            map[string]string{"doc.currency": "MXN"},
				privateKeyFile: "testdata/id_es256",
			},
			args: []string{"testdata/success.json"},
		},
		{
			name: "template missing",
			opts: &buildOpts{
				template: "missing.yaml",
			},
			err: "open missing.yaml: no such file or directory",
		},
		{
			name: "template",
			in:   strings.NewReader("{}"),
			opts: &buildOpts{
				template:       "testdata/success.yaml",
				privateKeyFile: "testdata/id_es256",
			},
		},
		{
			name: "type on command line",
			in:   readFile(t, "testdata/notype.json"),
			opts: &buildOpts{
				privateKeyFile: "testdata/id_es256",
				docType:        "bill.Invoice",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &cobra.Command{}
			if tt.in != nil {
				c.SetIn(tt.in)
			}
			buf := &bytes.Buffer{}
			c.SetOut(buf)
			opts := tt.opts
			if opts == nil {
				opts = &buildOpts{}
			}
			err := opts.runE(c, tt.args)
			if tt.err != "" {
				assert.EqualError(t, err, tt.err)
			} else {
				assert.Nil(t, err)
			}
			re := testy.Replacement{
				Regexp:      regexp.MustCompile(`(?sm)"sigs":.?\[.*\]`),
				Replacement: `"sigs": ["sig data"]`,
			}

			if d := testy.DiffText(testy.Snapshot(t), buf.String(), re); d != nil {
				t.Error(d)
			}
			if tt.target != "" {
				result, err := ioutil.ReadFile(tt.target)
				if err != nil {
					t.Fatal(err)
				}
				if d := testy.DiffText(testy.Snapshot(t, "outfile"), result, re); d != nil {
					t.Error(d)
				}
			}
		})
	}
}