/*
Copyright © 2019 Evan Cordell <cordell.evan@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ecordell/kpg/pkg/bundle"
	"github.com/ecordell/kpg/pkg/signals"
)

type pullOptions struct {
	outputDir string

	configs  []string
	username string
	password string
}

var pullOpts pullOptions

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := signals.Context()

		if len(args) < 1  {
			return fmt.Errorf("should be called with  one arg: host")
		}
		host := args[0]
		resolver := newResolver(pushOpts.username, pushOpts.password, pushOpts.configs...)
		return bundle.Pull(ctx, resolver, host, pullOpts.outputDir)
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)

	pullCmd.Flags().StringVarP(&pullOpts.outputDir, "out", "o", "", "directory to place files")
	pullCmd.Flags().StringArrayVarP(&pullOpts.configs, "config", "c", []string{"~/.docker/config.json"}, "auth config path")
	pullCmd.Flags().StringVarP(&pullOpts.username, "username", "u", "", "registry username")
	pullCmd.Flags().StringVarP(&pullOpts.password, "password", "p", "", "registry password")
}
