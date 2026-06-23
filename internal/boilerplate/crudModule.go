package boilerplate

import (
	"fmt"
	"path"
)

type CRUDModuleBoilerplate struct {
	CommonModuleBoilerplate
}

func (b *CRUDModuleBoilerplate) Create() (string, error) {
	tmplTypeName := "crud"
	err := b.CommonCreate(tmplTypeName)
	if err != nil {
		return "", err
	}

	// Create REST docs
	err = b.CreateRESTDocDir()
	if err != nil {
		return "", err
	}

	tmplData := NewModuleTmplData(b.Name)

	schemasTmplPath := path.Join(b.GetModuleTemplatesPath(), tmplTypeName, "schemas.yml.tmpl")
	schemasFilePath := path.Join(b.GetModuleRESTDocPath(), "schemas.yml")
	err = b.CreateFileFromTemplate(schemasTmplPath, schemasFilePath, tmplData)
	if err != nil {
		return "", err
	}

	resourceTmplPath := path.Join(b.GetModuleTemplatesPath(), tmplTypeName, "resource.yml.tmpl")
	resourceFilePath := path.Join(b.GetModuleRESTDocPath(), fmt.Sprintf("%s.yml", b.Name))
	err = b.CreateFileFromTemplate(resourceTmplPath, resourceFilePath, tmplData)
	if err != nil {
		return "", err
	}

	resourceIDTmplPath := path.Join(b.GetModuleTemplatesPath(), tmplTypeName, "resource-id.yml.tmpl")
	resourceIDFilePath := path.Join(b.GetModuleRESTDocPath(), fmt.Sprintf("%s-id.yml", b.Name))
	err = b.CreateFileFromTemplate(resourceIDTmplPath, resourceIDFilePath, tmplData)
	if err != nil {
		return "", err
	}

	openAPIPartPath := path.Join(b.GetModuleTemplatesPath(), tmplTypeName, "openapi-part.yml.tmpl")
	openAPIPartRes, err := b.RenderFromTemplate(openAPIPartPath, tmplData)
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf(
		"Created.\nPlease add strings below to the openapi.yml "+
			"for the REST API doc:\n\n%s",
		openAPIPartRes,
	)

	return result, nil
}
