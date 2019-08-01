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
	"context"
	"fmt"
	"os"

	"github.com/containerd/containerd/remotes"
	"github.com/containerd/containerd/remotes/docker"
	auth "github.com/deislabs/oras/pkg/auth/docker"
	"github.com/spf13/cobra"

	"github.com/ecordell/kpg/pkg/bundle"
	"github.com/ecordell/kpg/pkg/signals"
)

type pushOptions struct {
	configs  []string
	username string
	password string
}

var pushOpts pushOptions

func newResolver(username, password string, configs ...string) remotes.Resolver {
	if username != "" || password != "" {
		return docker.NewResolver(docker.ResolverOptions{
			Credentials: func(hostName string) (string, string, error) {
				return username, password, nil
			},
		})
	}
	cli, err := auth.NewClient(configs...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: Error loading auth file: %v\n", err)
	}
	resolver, err := cli.Resolver(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: Error loading resolver: %v\n", err)
		resolver = docker.NewResolver(docker.ResolverOptions{})
	}
	return resolver
}

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := signals.Context()

		if len(args) < 2  {
			return fmt.Errorf("should be called with two args: dir host")
		}
		dir := args[0]
		host := args[1]

		ctx = context.WithValue(ctx, bundle.BlobReaderContextKey, bundle.KustomizeBaseBlobReader)
		b, err := bundle.Build(ctx, dir)
		if err != nil {
			return err
		}
		resolver := newResolver(pushOpts.username, pushOpts.password, pushOpts.configs...)
		return bundle.Push(ctx, resolver, host, b)
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().StringArrayVarP(&pushOpts.configs, "config", "c", []string{"~/.docker/config.json"}, "auth config path")
	pushCmd.Flags().StringVarP(&pushOpts.username, "username", "u", "", "registry username")
	pushCmd.Flags().StringVarP(&pushOpts.password, "password", "p", "", "registry password")
}
