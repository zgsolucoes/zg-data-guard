package main

import (
	"github.com/zgsolucoes/zg-data-guard/config"
	"github.com/zgsolucoes/zg-data-guard/internal/webserver/router"
)

// @title           ZG Data Guard API
// @version         1.0
// @description     <b>This is the ZG Data Guard API.</b> <br />
// @description     ZG Data Guard is a centralized tool designed to streamline and secure the management of multiple databases across various environments.
// @description     It simplifies administration by providing a unified platform to handle database ecosystems, technologies, instances, predefined roles,
// @description     databases, users, and access control. All secured through JWT-protected API. <br /><br />Enjoy :D
// @termsOfService  http://swagger.io/terms/

// @contact.name   Luiz Henrique F. da Silva
// @contact.email  luizhenrique@zgsolucoes.com.br

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @Tag.name Ecosystem
// @Tag.description It represents the ecosystem where the database instance (cluster) is running. e.g. AWS, Cloud XPTO, On-premises, etc.
// @Tag.name Technology
// @Tag.description It represents the database technology. e.g. PostgreSQL 13, PostgreSQL 16, MySQL 5, etc.
// @Tag.name Database Instance
// @Tag.description It represents the database instance (cluster) that is running in a specific ecosystem.
// @Tag.name Database
// @Tag.description It represents the database that is running in a specific database instance. e.g. zg-data-guard, users-service, foo-service, etc.
// @Tag.name Database Role
// @Tag.description It represents the role that can be assigned to a user in a specific database. Each role has specific permissions. e.g. user_ro, developer, devops, etc.
// @Tag.name Database User
// @Tag.description It represents the user that can be created in a specific database instance (cluster) with a specific role. e.g. foo.bar, john.doe, etc.
// @Tag.name Access Permission
// @Tag.description It represents the permission that can be granted to a user to connect in a specific database. e.g. foo.bar (user) can connect in zg-data-guard (database) with developer role.

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Initialize Configs: Envs, Database Connection, Migrations, JWT, Crypto, etc
	config.Init()
	// Initialize WebServer
	router.Init()
	config.Cleanup()
}
