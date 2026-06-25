package app

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const pkgNamePlaceholder = "{{ .PkgName }}"

//go:embed templates/*
var templatesFiles embed.FS

type AppGenerator interface {
	CreateTemplate(boilerplateDir, templateDir, pkgName string) (err error)
	Create(templateName, appDir, pkgName string) (err error)
}

type AppGeneratorBasic struct {
}

func NewAppGenerator() (AppGenerator, error) {
	return &AppGeneratorBasic{}, nil
}

func (b *AppGeneratorBasic) CreateTemplate(boilerplateDir, templateDir, pkgName string) error {
	return b.Transform(
		nil,
		boilerplateDir,
		templateDir,
		pkgName,
		pkgNamePlaceholder,
		func(path string, replaced bool) string {
			if replaced && !strings.HasSuffix(path, ".tmpl") {
				return path + ".tmpl"
			}
			return path
		},
	)
}

func (b *AppGeneratorBasic) Create(templateName, appDir, pkgName string) error {
	templatesFS, err := fs.Sub(templatesFiles, "templates")
	if err != nil {
		return err
	}

	templateDir := templateName

	return b.Transform(
		templatesFS,
		templateDir,
		appDir,
		pkgNamePlaceholder,
		pkgName,
		func(path string, _ bool) string {
			return strings.TrimSuffix(path, ".tmpl")
		},
	)
}

func (b *AppGeneratorBasic) Transform(
	sourceFS fs.FS,
	sourceDir string,
	destinationDir string,
	from string,
	to string,
	renameFn func(path string, replaced bool) string,
) error {
	if _, err := os.Stat(destinationDir); err == nil {
		return fmt.Errorf("destination already exists: %s", destinationDir)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if sourceFS == nil {
		sourceFS = os.DirFS(sourceDir)
		sourceDir = "."
	} else if sourceDir == "" {
		sourceDir = "."
	}

	tmpDir, err := os.MkdirTemp("", "transform-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	workspaceDir := filepath.Join(tmpDir, "workspace")

	if err := b.copyDir(sourceFS, sourceDir, workspaceDir); err != nil {
		return err
	}

	if err := b.transformTree(workspaceDir, from, to, renameFn); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(destinationDir), 0o755); err != nil {
		return err
	}

	if err := os.Rename(workspaceDir, destinationDir); err == nil {
		return nil
	}

	return b.copyDir(os.DirFS(workspaceDir), ".", destinationDir)
}

func (b *AppGeneratorBasic) transformTree(
	root string,
	from string,
	to string,
	renameFn func(path string, replaced bool) string,
) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		replaced, err := b.transformFile(path, from, to)
		if err != nil {
			return err
		}

		newPath := renameFn(path, replaced)

		if newPath != path {
			if err := os.Rename(path, newPath); err != nil {
				return err
			}
		}

		return nil
	})
}

func (b *AppGeneratorBasic) transformFile(path, from, to string) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}

	if !bytes.Contains(data, []byte(from)) {
		return false, nil
	}

	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	updated := bytes.ReplaceAll(
		data,
		[]byte(from),
		[]byte(to),
	)

	if err := os.WriteFile(path, updated, info.Mode()); err != nil {
		return false, err
	}

	return true, nil
}

func (b *AppGeneratorBasic) copyDir(sourceFS fs.FS, src, dst string) error {
	return fs.WalkDir(sourceFS, src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		target := filepath.Join(dst, rel)

		info, err := d.Info()
		if err != nil {
			return err
		}

		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		return b.copyFile(sourceFS, path, target, info.Mode())
	})
}

func (b *AppGeneratorBasic) copyFile(sourceFS fs.FS, src, dst string, mode fs.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	in, err := sourceFS.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	mode |= 0o200 // ensure owner write permission for later transformation
	out, err := os.OpenFile(
		dst,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		mode,
	)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
