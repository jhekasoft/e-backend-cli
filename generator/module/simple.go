package boilerplate

type SimpleModuleGenerator struct {
	CommonModuleGenerator
}

func (b *SimpleModuleGenerator) Create() (string, error) {
	tmplTypeName := "simple"
	err := b.CommonCreate(tmplTypeName)
	if err != nil {
		return "", err
	}

	return "Created.", nil
}
