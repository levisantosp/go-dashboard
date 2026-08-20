package status

import (
	"dash/plugins"

	"github.com/danielgtaylor/huma/v2"
)

func StatusRoutes(api huma.API) {
	group := huma.NewGroup(api, "/status")
	group.UseMiddleware(plugins.AuthPlugin(api))

	huma.Get(group, "/{id}", Get)
	huma.Get(group, "/count", Count)
}
