package boilerplate

type SimpleModuleBoilerplate struct {
	CommonModuleBoilerplate
}

func (b *SimpleModuleBoilerplate) Create() (string, error) {
	tmplTypeName := "simple"
	err := b.CommonCreate(tmplTypeName)
	if err != nil {
		return "", err
	}

	return "Created.", nil
}
