package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"gitlab.com/flimzy/testy"
)

func Test_verify(t *testing.T) {
	tests := []struct {
		name string
		in   io.Reader
		args []string
		err  string
	}{
		{
			name: "invalid stdin",
			in:   strings.NewReader("this isn't JSON"),
			err:  "code=400, message=error unmarshaling JSON: json: cannot unmarshal string into Go value of type gobl.Envelope",
		},
		{
			name: "read error",
			in:   testy.ErrorReader(`{"foo":`, errors.New("read error")),
			err:  "read error",
		},
		{
			name: "empty envelope",
			in:   strings.NewReader(`{}`),
			err:  "code=422, message=$schema: cannot be blank; doc: cannot be blank; head: cannot be blank.",
		},
		{
			name: "success",
			in: func() io.Reader {
				f, err := os.Open("testdata/signed.json")
				if err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					_ = f.Close()
				})
				return f
			}(),
		},
		{
			name: "digest mismatch",
			in: func() io.Reader {
				f, err := os.Open("testdata/digest-mismatch.json")
				if err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					_ = f.Close()
				})
				return f
			}(),
			err: "code=422, message=digest mismatch",
		},
		{
			name: "read from file",
			args: []string{"testdata/digest-mismatch.json"},
			err:  "code=422, message=digest mismatch",
		},
		{
			name: "file missing",
			args: []string{"asdf"},
			err:  "open asdf: no such file or directory",
		},
		{
			name: "explicit stdin",
			args: []string{"-"},
			in: func() io.Reader {
				f, err := os.Open("testdata/digest-mismatch.json")
				if err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					_ = f.Close()
				})
				return f
			}(),
			err: "code=422, message=digest mismatch",
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
			err := verify().runE(c, tt.args)
			if tt.err != "" {
				assert.EqualError(t, err, tt.err)
			} else {
				assert.Nil(t, err)
			}
			if d := testy.DiffText(testy.Snapshot(t), buf.String()); d != nil {
				t.Error(d)
			}
		})
	}
}
