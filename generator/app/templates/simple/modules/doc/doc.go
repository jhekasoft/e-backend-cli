package doc

import (
	"path"

	internalModels "github.com/jhekasoft/e-backend/models"
)

const DataPath = "./modules/doc/data"

type DocModule struct {
}

func (m *DocModule) Name() string {
	return "Doc"
}

func (m *DocModule) Run(c *internalModels.Core) error {
	c.Echo.Static("/doc", path.Join(DataPath, "public"))

	return nil
}

func NewModule() internalModels.Module {
	return &DocModule{}
}
