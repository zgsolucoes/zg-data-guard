package handler

import (
	"encoding/json"
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/config"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
)

func HealthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	healthcheck := &dto.HealthCheckOutputDTO{
		ServiceName:    config.GetAppName(),
		ServiceVersion: config.GetBuildInfo().Version,
		ContextPath:    config.GetAppContextPath(),
		BuildTime:      config.GetBuildInfo().BuildTime,
	}
	_ = json.NewEncoder(w).Encode(healthcheck)
}
