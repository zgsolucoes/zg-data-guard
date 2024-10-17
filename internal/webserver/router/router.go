package router

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/zgsolucoes/zg-data-guard/docs"

	"github.com/zgsolucoes/zg-data-guard/config"
)

const (
	rootPath     = "/"
	apiBasePath  = "/api"
	apiVersionV1 = "/v1"
)

func Init() {
	r := chi.NewRouter()
	// RealIP middleware will set the request's IP to the value of the X-Forwarded-For or X-Real-IP headers. Useful when server is behind a reverse proxy
	r.Use(middleware.RealIP)
	// Recover from panics without crashing server
	r.Use(middleware.Recoverer)
	// CleanPath middleware will clean up the request URL path, redirecting to the clean path. For example, /path//to will redirect to /path/to
	r.Use(middleware.CleanPath)

	basePath := config.GetAppContextPath()
	log.Printf("Application BasePath: %s", basePath)
	setupSwaggerInfo(basePath)

	// Configure Routes
	initializeRoutes(r, basePath)

	// Create the HTTP server
	webServer := &http.Server{
		Addr:              ":" + config.GetWebPort(),
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second, // 5 seconds
	}

	initializeServerWithGracefulShutdown(webServer)
}

func setupSwaggerInfo(basePath string) {
	if config.GetEnvironment() == config.EnvDevelopment {
		docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", config.GetExternalHost(), config.GetWebPort())
		docs.SwaggerInfo.BasePath = apiBasePath + apiVersionV1
		config.SetApplicationURL(fmt.Sprintf("http://%s:%s", config.GetExternalHost(), config.GetWebPort()))
	} else {
		docs.SwaggerInfo.Host = config.GetExternalHost()
		docs.SwaggerInfo.BasePath = basePath + apiBasePath + apiVersionV1
		config.SetApplicationURL(fmt.Sprintf("https://%s%s", config.GetExternalHost(), basePath))
	}
}

func initializeServerWithGracefulShutdown(webServer *http.Server) {
	// Channel to receive OS signals
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Run the server in a goroutine so that it doesn't block the main thread
	go func() {
		log.Printf("Server running on port %s. External URL: %s", config.GetWebPort(), config.GetApplicationURL())
		if err := webServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new requests.")
	}()

	// Wait for OS signals to shut down the server
	<-shutdownSignal
	log.Println("Interrupt signal received. Shutting down server...")
	ctx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := webServer.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	log.Println("Graceful shutdown completed on server.")
}
