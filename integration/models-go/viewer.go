package models

import "github.com/tinhtran24/gqlgen/integration/remote_api"

type Viewer struct {
	User *remote_api.User
}
