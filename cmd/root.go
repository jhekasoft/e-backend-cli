/*
Copyright © 2024 Eugene Efremov <jhekasoft@gmail.com>

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
	"os"
	"runtime/debug"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	colorSuccess = color.New(color.FgHiGreen)
	colorError   = color.New(color.FgHiRed)
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "e-backend-cli",
	Short: "e-backend-cli",
	Long:  banner(),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(printBanner)
}

func banner() string {
	cLogo := color.New(color.FgBlue)
	cLogo2 := color.New(color.FgYellow)
	banner := cLogo.Sprintf(`
▗▄▄▄▖▗▄▄▖  ▗▄▖  ▗▄▄▖▗▖ ▗▖▗▄▄▄▖▗▖  ▗▖▗▄▄▄ 
▐▌   ▐▌ ▐▌▐▌ ▐▌▐▌   ▐▌▗▞▘▐▌   ▐▛▚▖▐▌▐▌  █
▐▛▀▀▘▐▛▀▚▖▐▛▀▜▌▐▌   ▐▛▚▖ ▐▛▀▀▘▐▌ ▝▜▌▐▌  █
▐▙▄▄▖▐▙▄▞▘▐▌ ▐▌▝▚▄▄▖▐▌ ▐▌▐▙▄▄▖▐▌  ▐▌▐▙▄▄▀`) + cLogo2.Sprintf(" CLI\n") +
		cLogo2.Sprintf("Version: %s", getVersion())

	// Remove the first newline for better formatting
	return strings.Replace(banner, "\n", "", 1)
}

func printBanner() {
	fmt.Print(banner() + "\n\n")
}

func getVersion() (version string) {
	version = "unknown"

	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "" {
			version = info.Main.Version
		}
	}

	return
}

// checkErr prints the msg with the prefix 'Error:' and exits with error code 1.
// If the msg is nil, it does nothing.
func checkErr(msg interface{}) {
	if msg != nil {
		colorError.Fprintln(os.Stderr, "Error:", msg)
		os.Exit(1)
	}
}
