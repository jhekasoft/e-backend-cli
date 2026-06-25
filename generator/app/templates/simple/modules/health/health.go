package health

import (
	"net/http"

	"github.com/jhekasoft/e-backend/models"

	"github.com/labstack/echo/v4"
)

type HealthResponse struct {
	OK        bool
	Version   string
	BuildTime string
}

type HealthModule struct {
	echo *echo.Echo
}

func (m *HealthModule) Name() string {
	return "Health"
}

func (m *HealthModule) Run(core *models.Core) error {
	m.echo = core.Echo
	m.echo.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, HealthResponse{
			OK:        true,
			Version:   core.Version,
			BuildTime: core.BuildTime,
		})
	})

	return nil
}

func NewModule() models.Module {
	return &HealthModule{}
}
