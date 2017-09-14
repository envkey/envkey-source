// Copyright © 2017 Dane Schneider <dane@envkey.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"

	"github.com/envkey/envkey-fetch/fetch"
	"github.com/envkey/envkey-source/shell"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var cacheDir string
var noCache bool
var force bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use: `eval $(envkey-source ENVKEY [flags])
or just:
  eval $(envkey-source [flags])
if you have a .env file in the current directory that includes ENVKEY=...`,
	Short: "Sets shell environment variables with an ENVKEY",
	Run: func(cmd *cobra.Command, args []string) {
		opts := fetch.FetchOptions{ShouldCache: !noCache, CacheDir: cacheDir}
		if len(args) > 0 {
			fmt.Println(shell.Source(args[0], force, opts))
		} else {
			godotenv.Load()
			envkey := os.Getenv("ENVKEY")

			if envkey != "" {
				fmt.Println(shell.Source(envkey, force, opts))
			} else {
				cmd.Help()
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.Flags().BoolVarP(&force, "force", "f", false, "overwrite existing environment variables and/or other entries in .env file")
	RootCmd.Flags().BoolVar(&noCache, "no-cache", false, "do NOT cache encrypted config as a local backup")
	RootCmd.Flags().StringVar(&cacheDir, "cache-dir", "", "cache directory (default is $HOME/.envkey/cache)")
}
