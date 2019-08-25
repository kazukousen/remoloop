package api

const ResourceUsersMe Resource = "/1/users/me"

type Me struct {
	Nickname string `json:"nickname"`
}
