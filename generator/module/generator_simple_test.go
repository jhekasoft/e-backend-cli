package boilerplate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSimpleModuleGenerator_Success(t *testing.T) {
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

	gen, err := NewModuleGenerator(pkgName, modName, "simple", modulesPath, restDocPath)
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

	dirs := []string{"models", "repository", "service", "handler"}
	for _, d := range dirs {
		p := filepath.Join(moduleDir, d, d+".go")
		if _, err := os.Stat(p); os.IsNotExist(err) {
			t.Fatalf("expected file not created: %s", p)
		}
	}

	// Negative case: modulesPath is a file -> Create should fail
}

func TestSimpleModuleGenerator_FailWhenModulesPathIsFile(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	// create a file where modulesPath should be
	modulesPath := filepath.Join(tmp, "modules")
	if err := os.WriteFile(modulesPath, []byte("not a dir"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	restDocPath := filepath.Join(modulesPath, "doc")

	pkgName := "github.com/example/project"
	modName := "example"

	gen, err := NewModuleGenerator(pkgName, modName, "simple", modulesPath, restDocPath)
	if err != nil {
		t.Fatalf("NewModuleGenerator error: %v", err)
	}

	if _, err := gen.Create(); err == nil {
		t.Fatalf("expected Create to fail when modulesPath is a file, but it succeeded")
	}
}
