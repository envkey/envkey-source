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
var envFile string
var shouldCache bool
var shouldNotCache bool
var force bool
var printVersion bool
var pamCompatible bool
var dotEnvCompatible bool
var verboseOutput bool
var timeoutSeconds float64
var retries uint8
var retryBackoff float64

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

		opts := fetch.FetchOptions{cacheEnabled, cacheDir, "envkey-source", version.Version, verboseOutput, timeoutSeconds, retries, retryBackoff}
		if len(args) > 0 {
			fmt.Println(shell.Source(args[0], force, opts, pamCompatible, dotEnvCompatible))
		} else {
			godotenv.Load(envFile)
			envkey := os.Getenv("ENVKEY")
			fmt.Println(shell.Source(envkey, force, opts, pamCompatible, dotEnvCompatible))
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
	RootCmd.Flags().StringVar(&envFile, "env-file", ".env", "ENVKEY-containing env file name")
	RootCmd.Flags().BoolVar(&verboseOutput, "verbose", false, "print verbose output (default is false)")
	RootCmd.Flags().Float64Var(&timeoutSeconds, "timeout", 10.0, "timeout in seconds for http requests")
	RootCmd.Flags().Uint8Var(&retries, "retries", 3, "number of times to retry requests on failure")
	RootCmd.Flags().Float64Var(&retryBackoff, "retryBackoff", 1, "retry backoff factor: {retryBackoff} * (2 ^ {retries - 1})")
	RootCmd.Flags().BoolVar(&pamCompatible, "pam-compatible", false, "change output format to be compatible with /etc/environment on Linux")
	RootCmd.Flags().BoolVar(&dotEnvCompatible, "dot-env-compatible", false, "change output to .env format")

	// differences between bash syntax and the /etc/environment format, as parsed by PAM
	// (https://github.com/linux-pam/linux-pam/blob/master/modules/pam_env/pam_env.c#L194)
	// - one variable per line
	// - "export " prefix is allowed, and has no effect
	// - cannot quote the variable name
	// - can quote the variable value
	//   (but this has no effect - there are no special sequences that need to be escaped)
	// - embedded quotes in values are treated as any other character (so should not be escaped)
	// - embedded newlines in values will disappear
	//   (a single backslash "escapes" the newline for parsing purposes, but in the actual
	//   environment the newline will not appear)
}
