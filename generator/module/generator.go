package boilerplate

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

//go:embed templates/*
var templatesFiles embed.FS

type ModuleGenerator interface {
	Create() (result string, err error)
}

func NewModuleGenerator(pkgName, name, template, modulesPath, restDocPath string) (ModuleGenerator, error) {
	templatesFS, err := fs.Sub(templatesFiles, "templates")
	if err != nil {
		return nil, err
	}

	switch template {
	case "crud":
		return &CRUDModuleGenerator{
			CommonModuleGenerator{templatesFS, pkgName, name, modulesPath, restDocPath},
		}, nil
	default:
		return &SimpleModuleGenerator{
			CommonModuleGenerator{templatesFS, pkgName, name, modulesPath, restDocPath},
		}, nil
	}
}

type CommonModuleGenerator struct {
	TemplatesFS fs.FS
	PkgName     string
	Name        string
	ModulesPath string
	RESTDocPath string
}

func (g *CommonModuleGenerator) GetModulePath() string {
	return path.Join(g.ModulesPath, g.Name)
}

func (g *CommonModuleGenerator) GetModuleRESTDocPath() string {
	return path.Join(g.RESTDocPath, g.Name)
}

func (g *CommonModuleGenerator) CommonCreate(tmplTypeName string) error {
	// Create init file
	err := g.CreateInitFile()
	if err != nil {
		return err
	}

	// Create module directory
	err = g.CreateModuleDir()
	if err != nil {
		return err
	}

	// Create module file
	err = g.CreateModuleFile(tmplTypeName)
	if err != nil {
		return err
	}

	// Create directories and files from templates in the module
	dirs := []string{"models", "repository", "service", "handler"}
	for _, dir := range dirs {
		err = g.CreateInModuleDir(dir)
		if err != nil {
			return err
		}

		tmplPath := path.Join(tmplTypeName, fmt.Sprintf("%s.go.tmpl", dir))
		filePath := path.Join(g.GetModulePath(), dir, fmt.Sprintf("%s.go", dir))
		err = g.CreateFileFromTemplate(tmplPath, filePath, NewModuleTmplData(g.PkgName, g.Name))
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *CommonModuleGenerator) CreateInitFile() error {
	initTmplPath := "init.go.tmpl"
	initFilePath := path.Join(g.ModulesPath, fmt.Sprintf("%s.go", g.Name))
	data := InitModuleTmplData{PkgName: g.PkgName, MdlName: g.Name}

	return g.CreateFileFromTemplate(initTmplPath, initFilePath, data)
}

func (g *CommonModuleGenerator) CreateModuleDir() error {
	// Create module directory
	modulePath := g.GetModulePath()
	if _, err := os.Stat(modulePath); os.IsNotExist(err) {
		// create directory
		if err := os.Mkdir(modulePath, 0754); err != nil {
			return err
		}
	}

	return nil
}

func (g *CommonModuleGenerator) CreateInModuleDir(name string) error {
	// Create directory in the module
	inModulePath := path.Join(g.GetModulePath(), name)
	if _, err := os.Stat(inModulePath); os.IsNotExist(err) {
		// create directory
		if err := os.Mkdir(inModulePath, 0754); err != nil {
			return err
		}
	}

	return nil
}

func (g *CommonModuleGenerator) CreateModuleFile(tmplTypeName string) error {
	moduleTmplPath := path.Join(tmplTypeName, "module.go.tmpl")
	moduleFilePath := path.Join(g.GetModulePath(), fmt.Sprintf("%s.go", g.Name))

	return g.CreateFileFromTemplate(moduleTmplPath, moduleFilePath, NewModuleTmplData(g.PkgName, g.Name))
}

func (g *CommonModuleGenerator) CreateFileFromTemplate(templateFilePath, filePath string, data any) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	tmpl, err := template.ParseFS(g.TemplatesFS, templateFilePath)
	if err != nil {
		return err
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, data)
	if err != nil {
		return err
	}

	_, err = file.WriteString(buf.String())
	if err != nil {
		return err
	}

	return nil
}

func (g *CommonModuleGenerator) RenderFromTemplate(templateFilePath string, data any) (string, error) {
	tmpl, err := template.ParseFS(g.TemplatesFS, templateFilePath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

type InitModuleTmplData struct {
	PkgName, MdlName string
}

type ModuleTmplData struct {
	PkgName, MdlName, MdlNameCap string
}

func NewModuleTmplData(pkgName, name string) ModuleTmplData {
	capitalizedName := cases.Title(language.English, cases.Compact).String(name)
	return ModuleTmplData{
		PkgName:    pkgName,
		MdlName:    name,
		MdlNameCap: capitalizedName,
	}
}

func (g *CommonModuleGenerator) CreateRESTDocDir() error {
	// Create module directory
	restDocPath := g.GetModuleRESTDocPath()
	if _, err := os.Stat(restDocPath); os.IsNotExist(err) {
		// create directory
		if err := os.Mkdir(restDocPath, 0754); err != nil {
			return err
		}
	}

	return nil
}

func (g *CommonModuleGenerator) AppendRESTDocPaths(yamlPathsFragment string) error {
	openAPIPath := path.Join(g.RESTDocPath, "openapi.yml")

	data, err := os.ReadFile(openAPIPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")

	pathsIdx := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == "paths:" {
			pathsIdx = i
			break
		}
	}

	if pathsIdx == -1 {
		return fmt.Errorf("paths section not found")
	}

	insertIdx := len(lines)

	for i := pathsIdx + 1; i < len(lines); i++ {
		line := lines[i]

		if strings.TrimSpace(line) == "" {
			continue
		}

		// Found next top-level section.
		if !strings.HasPrefix(line, " ") &&
			!strings.HasPrefix(line, "\t") {
			insertIdx = i
			break
		}
	}

	fragmentLines := strings.Split(
		strings.TrimRight(yamlPathsFragment, "\n"),
		"\n",
	)

	newLines := append(
		lines[:insertIdx],
		append(fragmentLines, lines[insertIdx:]...)...,
	)

	return os.WriteFile(
		openAPIPath,
		[]byte(strings.Join(newLines, "\n")),
		0644,
	)
}
