// Copyright © 2017 Envkey Inc. <support@envkey.com>
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
	"github.com/envkey/envkey-source/version"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var cacheDir string
var shouldCache bool
var shouldNotCache bool
var force bool
var printVersion bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use: `eval $(envkey-source [flags])

You'll need .env file in the current directory that includes ENVKEY=... (in development) or an ENVKEY environment variable set (on a server).

You can also pass an ENVKEY directly (not recommended for real workflows):

  eval $(envkey-source ENVKEY [flags])`,

	Short: "Sets shell environment variables with an ENVKEY",
	Run: func(cmd *cobra.Command, args []string) {
		if printVersion {
			fmt.Println(version.Version)
			return
		}

		// Determine whether local caching for offline work should be enabled
		// yes if --cache flag or .env file (unless --no-cache flag)
		cacheEnabled := !shouldNotCache && shouldCache
		if !cacheEnabled && !shouldNotCache {
			if _, err := os.Stat(".env"); !os.IsNotExist(err) {
				cacheEnabled = true
			}
		}

		opts := fetch.FetchOptions{cacheEnabled, cacheDir}
		if len(args) > 0 {
			fmt.Println(shell.Source(args[0], force, opts))
		} else {
			godotenv.Load()
			envkey := os.Getenv("ENVKEY")
			fmt.Println(shell.Source(envkey, force, opts))
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
	RootCmd.Flags().BoolVarP(&printVersion, "version", "v", false, "prints the version")
	RootCmd.Flags().BoolVar(&shouldCache, "cache", false, "cache encrypted config as a local backup (default is true when .env file exists, false otherwise)")
	RootCmd.Flags().BoolVar(&shouldNotCache, "no-cache", false, "do NOT cache encrypted config as a local backup even when .env file exists")
	RootCmd.Flags().StringVar(&cacheDir, "cache-dir", "", "cache directory (default is $HOME/.envkey/cache)")
}
