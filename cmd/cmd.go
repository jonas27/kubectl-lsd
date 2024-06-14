package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/jonas27/kubectl-lsd/internal"
	"github.com/spf13/cobra"
)

func LSD(in []byte) (string, error) {
	return internal.Run(in)
}

// https://github.com/spf13/cobra/issues/1336#issuecomment-773598580
// Execute adds all child commands to the root command and
// executes the selected command.
func Execute() error {
	var inputFile *string

	rootCmd := &cobra.Command{
		Use: "lsd",
		Example: `kubectl get secret my-secret -o yaml | kubectl lsd
kubectl get secret mysecret -oyaml | kubectl lsd
kubectl lsd -f - <./my-secret.json
kubectl lsd -f ./my-secret.json
kubectl lsd -f ./my-secret.json`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var in []byte
			var err error
			if *inputFile == "-" {
				stdin := cmd.InOrStdin()
				if in, err = io.ReadAll(stdin); err != nil {
					return fmt.Errorf("error reading from stdin: %w", err)
				}
			} else {
				in, err = os.ReadFile(*inputFile)
				if err != nil {
					return fmt.Errorf("error reading file %s: %w", *inputFile, err)
				}
			}
			out, err := LSD(in)
			if err != nil {
				return fmt.Errorf("error running Lsd: %w", err)
			}
			cmd.Print(out)
			return nil
		},
	}

	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)

	inputFile = rootCmd.Flags().StringP("file", "f", "-", "file path to lsd, or - to read from stdin")
	if err := rootCmd.MarkFlagFilename("file"); err != nil {
		return fmt.Errorf("error marking flag filename: %w", err)
	}

	getCmd := getCmd()
	rootCmd.AddCommand(getCmd)

	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("error executing root command: %w", err)
	}
	return nil
}

func getCmd() *cobra.Command {
	kubectl := "kubectl"
	getCmd := &cobra.Command{
		Use: "get",
		Example: `kubectl lsd get -- secret mysecret -oyaml
kubectl lsd get -- secret -n default mysecret -ojson`,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true}, // do not validate kubectl get flags
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			kubectlCmd := exec.Command(kubectl, args...)
			kres, err := kubectlCmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("error invoking kubectl as %v: %w", kubectlCmd.Args, err)
			}

			out, err := LSD(kres)
			if err != nil {
				return fmt.Errorf("error running LSD: %w", err)
			}
			cmd.Println(out)
			return nil
		},
	}
	return getCmd
}
