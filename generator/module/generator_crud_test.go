package boilerplate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCRUDModuleGenerator_Success(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	modulesPath := filepath.Join(tmp, "modules")
	restDocPath := filepath.Join(modulesPath, "doc", "data", "public", "restapi", "openapi")

	if err := os.MkdirAll(modulesPath, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(restDocPath, 0755); err != nil {
		t.Fatal(err)
	}

	pkgName := "github.com/example/project"
	modName := "example"

	gen, err := NewModuleGenerator(pkgName, modName, "crud", modulesPath, restDocPath)
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
	if !strings.Contains(res, "Please add strings below to the openapi.yml") {
		t.Fatalf("expected openapi hint in result, got: %q", res)
	}

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

	b, err := os.ReadFile(resource)
	if err != nil {
		t.Fatalf("read resource file: %v", err)
	}
	capName := strings.ToUpper(modName[:1]) + modName[1:]
	if !strings.Contains(string(b), modName) && !strings.Contains(string(b), capName) {
		t.Fatalf("resource file does not contain module name: %s", resource)
	}
}

func TestCRUDModuleGenerator_FailWhenRESTDocPathIsFile(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	modulesPath := filepath.Join(tmp, "modules")
	if err := os.MkdirAll(modulesPath, 0755); err != nil {
		t.Fatal(err)
	}

	// create a file where restDocPath should be
	restDocPath := filepath.Join(modulesPath, "docfile")
	if err := os.WriteFile(restDocPath, []byte("not a dir"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	pkgName := "github.com/example/project"
	modName := "example"

	gen, err := NewModuleGenerator(pkgName, modName, "crud", modulesPath, restDocPath)
	if err != nil {
		t.Fatalf("NewModuleGenerator error: %v", err)
	}

	if _, err := gen.Create(); err == nil {
		t.Fatalf("expected Create to fail when restDocPath parent is a file, but it succeeded")
	}
}
