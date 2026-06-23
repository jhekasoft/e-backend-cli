package boilerplate

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path"
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

func (b *CommonModuleGenerator) GetModulePath() string {
	return path.Join(b.ModulesPath, b.Name)
}

func (b *CommonModuleGenerator) GetModuleRESTDocPath() string {
	return path.Join(b.RESTDocPath, b.Name)
}

func (b *CommonModuleGenerator) CommonCreate(tmplTypeName string) error {
	// Create init file
	err := b.CreateInitFile()
	if err != nil {
		return err
	}

	// Create module directory
	err = b.CreateModuleDir()
	if err != nil {
		return err
	}

	// Create module file
	err = b.CreateModuleFile(tmplTypeName)
	if err != nil {
		return err
	}

	// Create directories and files from templates in the module
	dirs := []string{"models", "repository", "service", "handler"}
	for _, dir := range dirs {
		err = b.CreateInModuleDir(dir)
		if err != nil {
			return err
		}

		tmplPath := path.Join(tmplTypeName, fmt.Sprintf("%s.go.tmpl", dir))
		filePath := path.Join(b.GetModulePath(), dir, fmt.Sprintf("%s.go", dir))
		err = b.CreateFileFromTemplate(tmplPath, filePath, NewModuleTmplData(b.PkgName, b.Name))
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *CommonModuleGenerator) CreateInitFile() error {
	initTmplPath := "init.go.tmpl"
	initFilePath := path.Join(b.ModulesPath, fmt.Sprintf("%s.go", b.Name))
	data := InitModuleTmplData{PkgName: b.PkgName, MdlName: b.Name}

	return b.CreateFileFromTemplate(initTmplPath, initFilePath, data)
}

func (b *CommonModuleGenerator) CreateModuleDir() error {
	// Create module directory
	modulePath := b.GetModulePath()
	if _, err := os.Stat(modulePath); os.IsNotExist(err) {
		// create directory
		if err := os.Mkdir(modulePath, 0754); err != nil {
			return err
		}
	}

	return nil
}

func (b *CommonModuleGenerator) CreateInModuleDir(name string) error {
	// Create directory in the module
	inModulePath := path.Join(b.GetModulePath(), name)
	if _, err := os.Stat(inModulePath); os.IsNotExist(err) {
		// create directory
		if err := os.Mkdir(inModulePath, 0754); err != nil {
			return err
		}
	}

	return nil
}

func (b *CommonModuleGenerator) CreateModuleFile(tmplTypeName string) error {
	moduleTmplPath := path.Join(tmplTypeName, "module.go.tmpl")
	moduleFilePath := path.Join(b.GetModulePath(), fmt.Sprintf("%s.go", b.Name))

	return b.CreateFileFromTemplate(moduleTmplPath, moduleFilePath, NewModuleTmplData(b.PkgName, b.Name))
}

func (b *CommonModuleGenerator) CreateFileFromTemplate(templateFilePath, filePath string, data any) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	tmpl, err := template.ParseFS(b.TemplatesFS, templateFilePath)
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

func (b *CommonModuleGenerator) RenderFromTemplate(templateFilePath string, data any) (string, error) {
	tmpl, err := template.ParseFS(b.TemplatesFS, templateFilePath)
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

func (b *CommonModuleGenerator) CreateRESTDocDir() error {
	// Create module directory
	restDocPath := b.GetModuleRESTDocPath()
	if _, err := os.Stat(restDocPath); os.IsNotExist(err) {
		// create directory
		if err := os.Mkdir(restDocPath, 0754); err != nil {
			return err
		}
	}

	return nil
}
