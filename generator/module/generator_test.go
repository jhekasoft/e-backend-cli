package boilerplate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestModuleGenerators(t *testing.T) {
	cases := []struct {
		tmpl          string
		checkREST     bool
		expectSnippet string
	}{
		{tmpl: "simple", checkREST: false, expectSnippet: "Created."},
		{tmpl: "crud", checkREST: true, expectSnippet: "Please add strings below to the openapi.yml"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.tmpl, func(t *testing.T) {
			t.Parallel()

			tmp := t.TempDir()
			modulesPath := filepath.Join(tmp, "modules")
			restDocPath := filepath.Join(modulesPath, "doc", "data", "public", "restapi", "openapi")

			// ensure parent paths exist
			if err := os.MkdirAll(modulesPath, 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.MkdirAll(restDocPath, 0755); err != nil {
				t.Fatal(err)
			}

			pkgName := "github.com/example/project"
			modName := "example"

			gen, err := NewModuleGenerator(pkgName, modName, tc.tmpl, modulesPath, restDocPath)
			if err != nil {
				t.Fatalf("NewModuleGenerator error: %v", err)
			}

			res, err := gen.Create()
			if err != nil {
				t.Fatalf("Create error: %v", err)
			}

			if !strings.Contains(res, "Created.") {
				t.Fatalf("unexpected result: %q", res)
			}

			if tc.checkREST && !strings.Contains(res, tc.expectSnippet) {
				t.Fatalf("expected result to contain openapi hint, got: %q", res)
			}

			// Common checks: init file and module files/directories
			initFile := filepath.Join(modulesPath, modName+".go")
			if _, err := os.Stat(initFile); os.IsNotExist(err) {
				t.Fatalf("init file not created: %s", initFile)
			}

			moduleDir := filepath.Join(modulesPath, modName)
			if fi, err := os.Stat(moduleDir); err != nil || !fi.IsDir() {
				t.Fatalf("module dir not created: %s (err=%v)", moduleDir, err)
			}

			moduleFile := filepath.Join(moduleDir, modName+".go")
			if _, err := os.Stat(moduleFile); os.IsNotExist(err) {
				t.Fatalf("module file not created: %s", moduleFile)
			}

			// files in subdirs
			dirs := []string{"models", "repository", "service", "handler"}
			for _, d := range dirs {
				p := filepath.Join(moduleDir, d, d+".go")
				if _, err := os.Stat(p); os.IsNotExist(err) {
					t.Fatalf("expected file not created: %s", p)
				}
			}

			// read a generated file and assert it contains module name or capitalized variant
			data, err := os.ReadFile(moduleFile)
			if err != nil {
				t.Fatalf("read generated module file: %v", err)
			}
			content := string(data)
			capName := strings.ToUpper(modName[:1]) + modName[1:]
			if !(strings.Contains(content, modName) || strings.Contains(content, capName)) {
				t.Fatalf("generated module file does not contain module name (%s): %s", modName, moduleFile)
			}

			if tc.checkREST {
				restDocDir := filepath.Join(restDocPath, modName)
				if fi, err := os.Stat(restDocDir); err != nil || !fi.IsDir() {
					t.Fatalf("rest doc dir not created: %s (err=%v)", restDocDir, err)
				}

				schemas := filepath.Join(restDocDir, "schemas.yml")
				resource := filepath.Join(restDocDir, modName+".yml")
				resourceID := filepath.Join(restDocDir, modName+"-id.yml")

				for _, p := range []string{schemas, resource, resourceID} {
					if _, err := os.Stat(p); os.IsNotExist(err) {
						t.Fatalf("expected rest doc file not created: %s", p)
					}
				}

				// check resource contains module name
				b, err := os.ReadFile(resource)
				if err != nil {
					t.Fatalf("read resource file: %v", err)
				}
				if !strings.Contains(string(b), modName) && !strings.Contains(string(b), capName) {
					t.Fatalf("resource file does not contain module name: %s", resource)
				}
			}
		})
	}
}
