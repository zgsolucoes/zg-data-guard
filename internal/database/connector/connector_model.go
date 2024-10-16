package connector

import "github.com/zgsolucoes/zg-data-guard/internal/entity"

type Database struct {
	Name        string
	CurrentSize string
}

type DatabaseRole struct {
	Name entity.RoleName
}

type DatabaseUser struct {
	Username string
	Password string
	Role     string
}
