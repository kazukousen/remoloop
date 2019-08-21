package client

type resource string

const resourceDevices resource = "/1/devices"

type usersMe struct {
	Nickname string `json:"nickname"`
}

const resourceUsersMe resource = "/1/users/me"

type device struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
