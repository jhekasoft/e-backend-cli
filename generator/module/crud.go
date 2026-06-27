package boilerplate

import (
	"fmt"
	"path"
)

type CRUDModuleGenerator struct {
	CommonModuleGenerator
}

func (g *CRUDModuleGenerator) Create() (string, error) {
	tmplTypeName := "crud"
	err := g.CommonCreate(tmplTypeName)
	if err != nil {
		return "", err
	}

	result := "Created without REST docs."

	// Create REST docs
	if g.RESTDocPath != "" {
		err = g.createRESTDocs(tmplTypeName)
		if err != nil {
			return "", err
		}

		result = "Created with REST docs."
	}

	return result, nil
}

func (g *CRUDModuleGenerator) createRESTDocs(tmplTypeName string) error {
	err := g.CreateRESTDocDir()
	if err != nil {
		return err
	}

	tmplData := NewModuleTmplData(g.PkgName, g.Name)

	schemasTmplPath := path.Join(tmplTypeName, "schemas.yml.tmpl")
	schemasFilePath := path.Join(g.GetModuleRESTDocPath(), "schemas.yml")
	err = g.CreateFileFromTemplate(schemasTmplPath, schemasFilePath, tmplData)
	if err != nil {
		return err
	}

	resourceTmplPath := path.Join(tmplTypeName, "resource.yml.tmpl")
	resourceFilePath := path.Join(g.GetModuleRESTDocPath(), fmt.Sprintf("%s.yml", g.Name))
	err = g.CreateFileFromTemplate(resourceTmplPath, resourceFilePath, tmplData)
	if err != nil {
		return err
	}

	resourceIDTmplPath := path.Join(tmplTypeName, "resource-id.yml.tmpl")
	resourceIDFilePath := path.Join(g.GetModuleRESTDocPath(), fmt.Sprintf("%s-id.yml", g.Name))
	err = g.CreateFileFromTemplate(resourceIDTmplPath, resourceIDFilePath, tmplData)
	if err != nil {
		return err
	}

	openAPIPartPath := path.Join(tmplTypeName, "openapi-part.yml.tmpl")
	openAPIPartRes, err := g.RenderFromTemplate(openAPIPartPath, tmplData)
	if err != nil {
		return err
	}

	// Write new paths to the openapi.yml file
	err = g.AppendRESTDocPaths(openAPIPartRes)
	if err != nil {
		return err
	}

	return nil
}
