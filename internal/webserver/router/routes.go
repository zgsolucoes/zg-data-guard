package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/zgsolucoes/zg-data-guard/docs"

	"github.com/zgsolucoes/zg-data-guard/config"
	"github.com/zgsolucoes/zg-data-guard/internal/webserver/handler"
)

func initializeRoutes(r *chi.Mux, basePath string) {
	handler.InitializeAPIDependencies()

	r.Get(buildPath(basePath, "/"), handler.HomeHandler)
	setupHealthCheckRoutes(r, basePath)
	setupAuthRoutes(r, basePath)
	setupProtectedAPIRoutes(r, basePath)
	setupSwaggerRoute(r, basePath)
}

func setupHealthCheckRoutes(r *chi.Mux, basePath string) {
	r.Get("/healthcheck/info", handler.HealthCheckHandler)
	r.Get(buildPath(basePath, "/healthcheck/info"), handler.HealthCheckHandler)
}

func setupAuthRoutes(r *chi.Mux, basePath string) {
	if config.GetEnvironment() == config.EnvDevelopment {
		r.Get("/auth/internal", handler.InternalUserAuthHandler)
	}
}

func setupProtectedAPIRoutes(r *chi.Mux, basePath string) {
	apiRouter := chi.NewRouter()
	apiRouter.Use(middleware.Logger)
	apiRouter.Route("/", func(r chi.Router) {
		// Middleware to get the token from the request and set it in the context
		apiRouter.Use(jwtauth.Verifier(config.GetJwtHelper().Jwt))
		// Middleware to check if the token is valid
		apiRouter.Use(jwtauth.Authenticator)
		createEcosystemRoutes(apiRouter)
		createTechnologyRoutes(apiRouter)
		createDatabaseInstanceRoutes(apiRouter)
		createDatabaseRoutes(apiRouter)
		createDatabaseRoleRoutes(apiRouter)
		createDatabaseUserRoutes(apiRouter)
		createAccessPermissionRoutes(apiRouter)
	})

	r.Mount(buildPath(basePath, apiBasePath+apiVersionV1), apiRouter)
}

func setupSwaggerRoute(r *chi.Mux, basePath string) {
	r.Get(buildPath(basePath, "/docs/*"), httpSwagger.Handler(httpSwagger.URL(config.GetApplicationURL()+"/docs/doc.json")))
}

func createEcosystemRoutes(r chi.Router) {
	r.Route("/ecosystem", func(r chi.Router) {
		r.Post("/", handler.CreateEcosystemHandler)
		r.Put("/", handler.UpdateEcosystemHandler)
		r.Get("/", handler.GetEcosystemHandler)
		r.Delete("/", handler.DeleteEcosystemHandler)
	})
	r.Get("/ecosystems", handler.ListEcosystemsHandler)
}

func createTechnologyRoutes(r chi.Router) {
	r.Route("/technology", func(r chi.Router) {
		r.Post("/", handler.CreateTechnologyHandler)
		r.Put("/", handler.UpdateTechnologyHandler)
		r.Get("/", handler.GetTechnologyHandler)
		r.Delete("/", handler.DeleteTechnologyHandler)
	})
	r.Get("/technologies", handler.ListTechnologiesHandler)
}

func createDatabaseInstanceRoutes(r chi.Router) {
	r.Route("/database-instance", func(r chi.Router) {
		r.Post("/", handler.CreateDatabaseInstanceHandler)
		r.Put("/", handler.UpdateDatabaseInstanceHandler)
		r.Get("/", handler.GetDatabaseInstanceHandler)
		r.Get("/credentials", handler.GetDatabaseInstanceCredentialsHandler)
		r.Patch("/change-status", handler.ChangeStatusDatabaseInstanceHandler)
		r.Post("/test-connection", handler.TestConnectionHandler)
		r.Post("/propagate-roles", handler.PropagateRolesHandler)
		r.Post("/sync-databases", handler.SyncDatabasesHandler)
	})
	r.Get("/database-instances", handler.ListDatabaseInstancesHandler)
}

func createDatabaseRoutes(r chi.Router) {
	r.Route("/database", func(r chi.Router) {
		r.Get("/", handler.GetDatabaseHandler)
		r.Post("/setup-roles", handler.SetupRolesHandler)
	})
	r.Get("/databases", handler.ListDatabasesHandler)
}

func createDatabaseRoleRoutes(r chi.Router) {
	r.Get("/database-roles", handler.ListDatabaseRolesHandler)
}

func createDatabaseUserRoutes(r chi.Router) {
	r.Route("/database-user", func(r chi.Router) {
		r.Post("/", handler.CreateDatabaseUserHandler)
		r.Get("/", handler.GetDatabaseUserHandler)
		r.Put("/", handler.UpdateDatabaseUserHandler)
		r.Get("/credentials", handler.GetDatabaseUserCredentialsHandler)
		r.Patch("/change-status", handler.ChangeStatusDatabaseUserHandler)
	})
	r.Get("/database-users", handler.ListDatabaseUsersHandler)
}

func createAccessPermissionRoutes(r chi.Router) {
	r.Route("/access-permission", func(r chi.Router) {
		r.Post("/grant", handler.GrantAccessHandler)
		r.Post("/revoke", handler.RevokeAccessHandler)
		r.Get("/logs", handler.ListAccessPermissionLogsHandler)
	})
	r.Get("/access-permissions", handler.ListAccessPermissionsHandler)
}

func buildPath(basePath, path string) string {
	if config.GetEnvironment() == config.EnvDevelopment {
		return path
	}
	if path == rootPath {
		return basePath
	}
	return basePath + path
}
