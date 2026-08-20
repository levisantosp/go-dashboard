package auth

import "github.com/danielgtaylor/huma/v2"

func AuthRoutes(api huma.API) {
	group := huma.NewGroup(api, "/auth")
	huma.Get(group, "/discord", Auth)
	huma.Get(group, "/discord/callback", Callback)
	huma.Post(group, "/signout", SignOut)
	huma.Get(group, "/session", GetSession)
}
