/*
Copyright © 2025 Eugene Efremov <jhekasoft@gmail.com>

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
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"

	moduleGenerator "github.com/jhekasoft/e-backend-cli/generator/module"
	eCmd "github.com/jhekasoft/e-backend/cmd"
)

// moduleCreateCmd represents the moduleCreate command
var moduleCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new module",
	Long: `Create a new module.
For example:

e-backend-cli module create myFirstModule`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			eCmd.CheckErr(fmt.Errorf("\"module create\" needs a name for the module"))
		}

		name := args[0]
		template, _ := cmd.Flags().GetString("template")

		fmt.Printf("Creating module '%s' with template '%s'\n", name, template)

		// Read project package name (from go.mod)
		pkgName, err := determineProjectPackageName()
		eCmd.CheckErr(err)

		modulesPath := "modules"
		restDocPath := "modules/doc/data/public/restapi/openapi"
		mg, err := moduleGenerator.NewModuleGenerator(pkgName, name, template, modulesPath, restDocPath)
		eCmd.CheckErr(err)
		result, err := mg.Create()
		eCmd.CheckErr(err)

		eCmd.ColorSuccess.Println(result)
	},
}

func init() {
	moduleCmd.AddCommand(moduleCreateCmd)

	moduleCreateCmd.Flags().StringP("template", "t", "simple", `A module template to use.
Available templates: simple, crud.`,
	)
}

func determineProjectPackageName() (string, error) {
	content, err := os.ReadFile("go.mod")
	if err != nil {
		return "", err
	}

	// Parse the file structure
	ast, err := modfile.Parse("go.mod", content, nil)
	if err != nil {
		return "", err
	}

	// Access the module path statement directly
	if ast.Module != nil {
		return ast.Module.Mod.Path, nil
	}

	return "", errors.New("no module statement found in go.mod")
}
