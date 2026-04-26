package health

import "github.com/Ashishkapoor1469/Nestgo/common"

type HealthController struct{}

func NewHealthController() *HealthController { return &HealthController{} }

func (c *HealthController) Prefix() string { return "/health" }

func (c *HealthController) Routes() []common.Route {
	return []common.Route{
		{Method: "GET", Path: "/", Handler: c.Check, Summary: "Health check"},
	}
}

func (c *HealthController) Check(ctx *common.Context) error {
	return ctx.OK(map[string]any{
		"status":  "ok",
		"service": "basic-api",
	})
}
