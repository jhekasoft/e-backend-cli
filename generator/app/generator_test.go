package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateTemplate(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	boilerplateDir := filepath.Join(tmp, "boilerplate")
	templateDir := filepath.Join(tmp, "templates", "custom")

	if err := os.MkdirAll(filepath.Join(boilerplateDir, "cmd"), 0o755); err != nil {
		t.Fatalf("mkdir boilerplate dir: %v", err)
	}

	sourceFile := filepath.Join(boilerplateDir, "cmd", "root.go")
	sourceContents := []byte(`package cmd

import "example.com/demo/internal"
`)
	if err := os.WriteFile(sourceFile, sourceContents, 0o644); err != nil {
		t.Fatalf("write boilerplate file: %v", err)
	}

	gen, err := NewAppGenerator()
	if err != nil {
		t.Fatalf("NewAppGenerator error: %v", err)
	}

	pkgName := "example.com/demo"
	if err := gen.CreateTemplate(boilerplateDir, templateDir, pkgName); err != nil {
		t.Fatalf("CreateTemplate error: %v", err)
	}

	generatedFile := filepath.Join(templateDir, "cmd", "root.go.tmpl")
	data, err := os.ReadFile(generatedFile)
	if err != nil {
		t.Fatalf("read generated template: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "{{ .PkgName }}") {
		t.Fatalf("expected placeholder to be preserved in template output, got %q", content)
	}
	if _, err := os.Stat(generatedFile); err != nil {
		t.Fatalf("expected generated template file to exist: %v", err)
	}
}

func TestCreate(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	appDir := filepath.Join(tmp, "generated-app")

	gen, err := NewAppGenerator()
	if err != nil {
		t.Fatalf("NewAppGenerator error: %v", err)
	}

	pkgName := "example.com/demo"
	if err := gen.Create("simple", appDir, pkgName); err != nil {
		t.Fatalf("Create error: %v", err)
	}

	mainFile := filepath.Join(appDir, "main.go")
	data, err := os.ReadFile(mainFile)
	if err != nil {
		t.Fatalf("read generated main.go: %v", err)
	}

	mainContent := string(data)
	if !strings.Contains(mainContent, `import "`+pkgName+`/cmd"`) {
		t.Fatalf("expected generated main.go to import %q, got %q", pkgName+"/cmd", mainContent)
	}

	goModFile := filepath.Join(appDir, "go.mod")
	goModData, err := os.ReadFile(goModFile)
	if err != nil {
		t.Fatalf("read generated go.mod: %v", err)
	}

	if !strings.Contains(string(goModData), "module "+pkgName) {
		t.Fatalf("expected generated go.mod to contain module %q", pkgName)
	}

	if _, err := os.Stat(filepath.Join(appDir, "cmd", "root.go")); err != nil {
		t.Fatalf("expected generated cmd/root.go to exist: %v", err)
	}
}
