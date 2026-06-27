package boilerplate

type SimpleModuleGenerator struct {
	CommonModuleGenerator
}

func (g *SimpleModuleGenerator) Create() (string, error) {
	tmplTypeName := "simple"
	err := g.CommonCreate(tmplTypeName)
	if err != nil {
		return "", err
	}

	return "Created.", nil
}
