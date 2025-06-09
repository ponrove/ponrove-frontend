package config

import (
	"github.com/ponrove/configura"
	"github.com/ponrove/ponrove-frontend/pkg/webclient"
	"github.com/ponrove/ponrunner"
)

var serverConfigInstance *configura.ConfigImpl

// ServerConfig returns a singleton instance of the server configuration.
func New() configura.Config {
	if serverConfigInstance == nil {
		serverConfigInstance = configura.NewConfigImpl()
		configura.LoadEnvironment(serverConfigInstance, ponrunner.SERVER_OPENFEATURE_PROVIDER_NAME, "")
		configura.LoadEnvironment(serverConfigInstance, ponrunner.SERVER_OPENFEATURE_PROVIDER_URL, "")
		configura.LoadEnvironment(serverConfigInstance, ponrunner.SERVER_PORT, int64(8080))
		configura.LoadEnvironment(serverConfigInstance, ponrunner.SERVER_WRITE_TIMEOUT, int64(10))
		configura.LoadEnvironment(serverConfigInstance, ponrunner.SERVER_READ_TIMEOUT, int64(10))
		configura.LoadEnvironment(serverConfigInstance, ponrunner.SERVER_REQUEST_TIMEOUT, int64(30))
		configura.LoadEnvironment(serverConfigInstance, ponrunner.SERVER_SHUTDOWN_TIMEOUT, int64(10))
		configura.LoadEnvironment(serverConfigInstance, ponrunner.SERVER_LOG_LEVEL, "info")
		configura.LoadEnvironment(serverConfigInstance, ponrunner.SERVER_LOG_FORMAT, "json")
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_ENABLED, false)
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_LOGS_ENABLED, false)
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_METRICS_ENABLED, false)
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_TRACES_ENABLED, false)
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_SERVICE_NAME, "ponrove-frontend")
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_EXPORTER_OTLP_ENDPOINT, "")
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_EXPORTER_OTLP_TRACES_ENDPOINT, "")
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_EXPORTER_OTLP_METRICS_ENDPOINT, "")
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_EXPORTER_OTLP_LOGS_ENDPOINT, "")
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_EXPORTER_OTLP_HEADERS, "")
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_EXPORTER_OTLP_TRACES_HEADERS, "")
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_EXPORTER_OTLP_METRICS_HEADERS, "")
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_EXPORTER_OTLP_LOGS_HEADERS, "")
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_EXPORTER_OTLP_TIMEOUT, int64(10))
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_EXPORTER_OTLP_TRACES_TIMEOUT, int64(10))
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_EXPORTER_OTLP_METRICS_TIMEOUT, int64(10))
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_EXPORTER_OTLP_LOGS_TIMEOUT, int64(10))
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_EXPORTER_OTLP_PROTOCOL, "grpc")
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_EXPORTER_OTLP_TRACES_PROTOCOL, "grpc")
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_EXPORTER_OTLP_METRICS_PROTOCOL, "grpc")
		configura.LoadEnvironment(serverConfigInstance, ponrunner.OTEL_EXPORTER_OTLP_LOGS_PROTOCOL, "grpc")
		configura.LoadEnvironment(serverConfigInstance, webclient.WEBCLIENT_APP_BUILD_DIR, "./app/build")
	}

	return *serverConfigInstance
}
