package boilerplate

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const TemplatesPath = "internal/boilerplate/templates"

type ModuleBoilerplate interface {
	Create() (result string, err error)
}

func NewModuleBoilerplate(name, template, modulesPath, restDocPath string) ModuleBoilerplate {
	switch template {
	case "crud":
		return &CRUDModuleBoilerplate{CommonModuleBoilerplate{name, modulesPath, restDocPath}}
	default:
		return &SimpleModuleBoilerplate{CommonModuleBoilerplate{name, modulesPath, restDocPath}}
	}
}

type CommonModuleBoilerplate struct {
	Name        string
	ModulesPath string
	RESTDocPath string
}

func (b *CommonModuleBoilerplate) GetModuleTemplatesPath() string {
	exePath, err := os.Executable()
	if err != nil {
		// TODO: move logging upper
		log.Fatalf("Error fetching executable path: %v", err)
	}
	return path.Join(filepath.Dir(exePath), TemplatesPath, "module")
}

func (b *CommonModuleBoilerplate) GetModulePath() string {
	return path.Join(b.ModulesPath, b.Name)
}

func (b *CommonModuleBoilerplate) GetModuleRESTDocPath() string {
	return path.Join(b.RESTDocPath, b.Name)
}

func (b *CommonModuleBoilerplate) CommonCreate(tmplTypeName string) error {
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

		tmplPath := path.Join(b.GetModuleTemplatesPath(), tmplTypeName, fmt.Sprintf("%s.go.tmpl", dir))
		filePath := path.Join(b.GetModulePath(), dir, fmt.Sprintf("%s.go", dir))
		err = b.CreateFileFromTemplate(tmplPath, filePath, NewModuleTmplData(b.Name))
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *CommonModuleBoilerplate) CreateInitFile() error {
	initTmplPath := path.Join(b.GetModuleTemplatesPath(), "init.go.tmpl")
	initFilePath := path.Join(b.ModulesPath, fmt.Sprintf("%s.go", b.Name))
	data := InitModuleTmplData{MdlName: b.Name}

	return b.CreateFileFromTemplate(initTmplPath, initFilePath, data)
}

func (b *CommonModuleBoilerplate) CreateModuleDir() error {
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

func (b *CommonModuleBoilerplate) CreateInModuleDir(name string) error {
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

func (b *CommonModuleBoilerplate) CreateModuleFile(tmplTypeName string) error {
	moduleTmplPath := path.Join(b.GetModuleTemplatesPath(), tmplTypeName, "module.go.tmpl")
	moduleFilePath := path.Join(b.GetModulePath(), fmt.Sprintf("%s.go", b.Name))

	return b.CreateFileFromTemplate(moduleTmplPath, moduleFilePath, NewModuleTmplData(b.Name))
}

func (b *CommonModuleBoilerplate) CreateFileFromTemplate(templateFilePath, filePath string, data any) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	tmpl, err := template.ParseFiles(templateFilePath)
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

func (b *CommonModuleBoilerplate) RenderFromTemplate(templateFilePath string, data any) (string, error) {
	tmpl, err := template.ParseFiles(templateFilePath)
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
	MdlName string
}

type ModuleTmplData struct {
	MdlName    string
	MdlNameCap string
}

func NewModuleTmplData(name string) ModuleTmplData {
	capitalizedName := cases.Title(language.English, cases.Compact).String(name)
	return ModuleTmplData{
		MdlName:    name,
		MdlNameCap: capitalizedName,
	}
}

func (b *CommonModuleBoilerplate) CreateRESTDocDir() error {
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
