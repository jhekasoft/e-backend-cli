/*
Copyright © 2026 Eugene Efremov <jhekasoft@gmail.com>

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

	appGenerator "github.com/jhekasoft/e-backend-cli/generator/app"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new application",
	Long: `Create a new application from template.
For example:

e-backend-cli app create --template simple --appDir /path/to/new/app --pkgName myapp`,
	Run: func(cmd *cobra.Command, args []string) {
		templateName, _ := cmd.Flags().GetString("template")
		appDir, _ := cmd.Flags().GetString("appDir")
		pkgName, _ := cmd.Flags().GetString("pkgName")

		appTemplateGenerator, err := appGenerator.NewAppGenerator()
		cobra.CheckErr(err)

		err = appTemplateGenerator.Create(templateName, appDir, pkgName)
		cobra.CheckErr(err)

		fmt.Printf("Application created successfully at: %s\n", appDir)
	},
}

func init() {
	appCmd.AddCommand(createCmd)

	createCmd.Flags().StringP("template", "t", "", `Template name for the new application (e.g., "simple")`)
	createCmd.Flags().StringP("appDir", "a", "", "Path to the directory of the new application")
	createCmd.Flags().StringP("pkgName", "p", "", `Application package name for replacement in the template (e.g., "myapp")`)

	createCmd.MarkFlagRequired("template")
	createCmd.MarkFlagRequired("appDir")
	createCmd.MarkFlagRequired("pkgName")
}
