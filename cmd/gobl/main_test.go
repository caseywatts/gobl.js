package main

import (
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/flimzy/testy"

	"github.com/invopop/gobl"
)

func Test_root(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		stdin io.Reader
		err   string
	}{
		{
			name: "unsupported command",
			args: []string{"foo"},
			err:  `unknown command "foo" for "gobl"`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cmd := root().cmd()
			cmd.SetArgs(tt.args)
			var err error
			stdout, stderr := testy.RedirIO(tt.stdin, func() {
				err = cmd.Execute()
			})
			if d := testy.DiffText(testy.Snapshot(t, "_stdout"), stdout); d != nil {
				t.Errorf("STDOUT: %s", d)
			}
			if d := testy.DiffText(testy.Snapshot(t, "_stderr"), stderr); d != nil {
				t.Errorf("STDERR: %s", d)
			}
			if tt.err == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, tt.err)
			}
		})
	}
}

func Test_version(t *testing.T) {
	cmd := version()
	stdout, stderr := testy.RedirIO(nil, func() {
		err := cmd.Execute()
		if err != nil {
			t.Fatal(err)
		}
	})
	wantOut := "GOBL version " + string(gobl.VERSION) + "\n"
	wantErr := ""
	if sout, _ := ioutil.ReadAll(stdout); string(sout) != wantOut {
		t.Errorf("Unexpected STDOUT: %s", sout)
	}
	if serr, _ := ioutil.ReadAll(stderr); string(serr) != wantErr {
		t.Errorf("Unexpected STDERR: %s", serr)
	}
}
