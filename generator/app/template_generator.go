package app

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const pkgNamePlaceholder = "{{ .PkgName }}"

type AppTemplateGenerator interface {
	Create(boilerplateDir, templateDir, pkgName string) (err error)
}

type AppTemplateGeneratorBasic struct {
}

func NewAppTemplateGenerator() (AppTemplateGenerator, error) {
	return &AppTemplateGeneratorBasic{}, nil
}

func (b *AppTemplateGeneratorBasic) Create(boilerplateDir, templateDir, pkgName string) (err error) {
	if _, err := os.Stat(templateDir); err == nil {
		return fmt.Errorf("template directory already exists: %s", templateDir)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	tmpDir, err := os.MkdirTemp("", "project-generator-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	workspaceDir := filepath.Join(tmpDir, "workspace")

	if err := b.copyDir(boilerplateDir, workspaceDir); err != nil {
		return fmt.Errorf("copy template: %w", err)
	}

	if err := b.processTree(workspaceDir, pkgName); err != nil {
		return fmt.Errorf("process template: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(templateDir), 0o755); err != nil {
		return fmt.Errorf("create destination parent: %w", err)
	}

	// Fast path: atomic move on the same filesystem.
	if err := os.Rename(workspaceDir, templateDir); err == nil {
		return nil
	}

	// Fallback: copy recursively.
	if err := b.copyDir(workspaceDir, templateDir); err != nil {
		return fmt.Errorf("copy generated project: %w", err)
	}

	return nil
}

func (b *AppTemplateGeneratorBasic) processTree(root, pkgName string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		shouldRename, err := b.processFile(path, pkgName)
		if err != nil {
			return err
		}

		if shouldRename && !strings.HasSuffix(path, ".tmpl") {
			if err := os.Rename(path, path+".tmpl"); err != nil {
				return fmt.Errorf("rename %s: %w", path, err)
			}
		}

		return nil
	})
}

func (b *AppTemplateGeneratorBasic) processFile(path, pkgName string) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}

	if !bytes.Contains(data, []byte(pkgName)) {
		return false, nil
	}

	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	updated := bytes.ReplaceAll(
		data,
		[]byte(pkgName),
		[]byte(pkgNamePlaceholder),
	)

	if err := os.WriteFile(path, updated, info.Mode()); err != nil {
		return false, err
	}

	return true, nil
}

func (b *AppTemplateGeneratorBasic) copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
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
			return os.MkdirAll(target, info.Mode())
		}

		return copyFile(path, target, info.Mode())
	})
}

func copyFile(src, dst string, mode fs.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(
		dst,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		mode,
	)
	if err != nil {
		return err
	}

	defer func() {
		_ = out.Close()
	}()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return out.Close()
}
