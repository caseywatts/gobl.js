// The gobl command provides a command-line interface to the GOBL library.

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl.cli/internal"
)

func main() {
	if err := run(); err != nil {
		echoErr := new(echo.HTTPError)
		if errors.As(err, &echoErr) {
			msg := echoErr.Message
			int := echoErr.Internal
			switch {
			case msg != "" && int != nil:
				err = fmt.Errorf("%v: %w", msg, int)
			case int != nil:
				err = int
			default:
				err = fmt.Errorf("%v", msg)
			}
		}
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	return root().ExecuteContext(ctx)
}

func root() *cobra.Command {
	root := &cobra.Command{
		Use:           "gobl",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(&cobra.Command{
		Use:  "verify [infile]",
		Args: cobra.MaximumNArgs(1),
		RunE: verify,
	})
	root.AddCommand(build().cmd())
	root.AddCommand(version())
	root.AddCommand(serve().cmd())
	root.AddCommand(keygen().cmd())
	return root
}

func inputFilename(args []string) string {
	if len(args) > 0 && args[0] != "-" {
		return args[0]
	}
	return ""
}

func openInput(cmd *cobra.Command, args []string) (io.ReadCloser, error) {
	if inFile := inputFilename(args); inFile != "" {
		return os.Open(inFile)
	}
	return ioutil.NopCloser(cmd.InOrStdin()), nil
}

func verify(cmd *cobra.Command, args []string) error {
	input, err := openInput(cmd, args)
	if err != nil {
		return err
	}
	defer input.Close() // nolint:errcheck

	return internal.Verify(cmdContext(cmd), input)
}

func version() *cobra.Command {
	return &cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, _ []string) {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "GOBL version %s\n", gobl.VERSION)
		},
	}
}
