// Command hello is the sample CLI that ships with the docsmith template. It
// has a greeting and a word counter, just enough for every generator to have
// something real to document.
package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/viniciusdc/docsmith/internal/textstat"
)

// version is overridden at build time with -ldflags "-X main.version=…".
var version = "dev"

func main() {
	if err := newRootCmd(os.Stdin, os.Stdout).Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd(stdin io.Reader, stdout io.Writer) *cobra.Command {
	var upper bool

	root := &cobra.Command{
		Use:   "hello",
		Short: "A tiny CLI that ships with the docsmith template",
		Long: `hello is the sample command-line tool in the docsmith template. Replace it
with your own; the generators in internal/tools/ only need a cobra binary.`,
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: false,
	}
	root.SetIn(stdin)
	root.SetOut(stdout)
	root.PersistentFlags().BoolVarP(&upper, "upper", "u", false, "print output in upper case")

	// out writes a line to stdout, upper-cased when --upper is set.
	out := func(cmd *cobra.Command, format string, args ...any) {
		line := fmt.Sprintf(format, args...)
		if upper {
			line = strings.ToUpper(line)
		}
		fmt.Fprintln(cmd.OutOrStdout(), line)
	}

	root.AddCommand(newGreetCmd(out), newCountCmd(out))
	return root
}

type printer func(cmd *cobra.Command, format string, args ...any)

func newGreetCmd(out printer) *cobra.Command {
	var greeting string
	cmd := &cobra.Command{
		Use:     "greet [name]",
		Short:   "Print a greeting",
		Example: "  hello greet\n  hello greet Ada --greeting Hi",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := "world"
			if len(args) == 1 {
				name = args[0]
			}
			out(cmd, "%s, %s!", greeting, name)
			return nil
		},
	}
	cmd.Flags().StringVarP(&greeting, "greeting", "g", "Hello", "word to greet with")
	return cmd
}

func newCountCmd(out printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "count <file|->...",
		Short:   "Count lines, words and bytes, like wc",
		Example: "  hello count notes.txt\n  cat notes.txt | hello count -",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var total textstat.Stats
			for _, name := range args {
				s, err := countOne(cmd, name)
				if err != nil {
					return err
				}
				out(cmd, "%6d %6d %6d %s", s.Lines, s.Words, s.Bytes, name)
				total = total.Add(s)
			}
			if len(args) > 1 {
				out(cmd, "%6d %6d %6d total", total.Lines, total.Words, total.Bytes)
			}
			return nil
		},
	}
	return cmd
}

func countOne(cmd *cobra.Command, name string) (textstat.Stats, error) {
	if name == "-" {
		return textstat.Count(cmd.InOrStdin())
	}
	f, err := os.Open(name)
	if err != nil {
		return textstat.Stats{}, err
	}
	defer func() { _ = f.Close() }()
	return textstat.Count(f)
}
