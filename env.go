package autolemetry

import (
	"net/url"
	"os"
	"strings"
)

func applyEnvOverrides(cfg *Config) {
	if cfg.ServiceName == "" || cfg.ServiceName == defaultServiceName {
		if v := strings.TrimSpace(os.Getenv("OTEL_SERVICE_NAME")); v != "" {
			cfg.ServiceName = v
		}
	}

	if cfg.Endpoint == "" {
		if v := strings.TrimSpace(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")); v != "" {
			cfg.Endpoint = sanitizeEndpoint(v)
		}
	}

	if proto := strings.ToLower(strings.TrimSpace(os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL"))); proto != "" {
		switch proto {
		case string(ProtocolHTTP):
			cfg.Protocol = ProtocolHTTP
		case string(ProtocolGRPC):
			cfg.Protocol = ProtocolGRPC
		}
	}

	if len(cfg.OTLPHeaders) == 0 {
		if headers := parseHeadersEnv(os.Getenv("OTEL_EXPORTER_OTLP_HEADERS")); len(headers) > 0 {
			cfg.OTLPHeaders = headers
		}
	}
}

func sanitizeEndpoint(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}

	// Attempt full URL parsing first; fall back to manual cleanup.
	if strings.Contains(raw, "://") {
		if u, err := url.Parse(raw); err == nil {
			if host := u.Host; host != "" {
				return strings.TrimSuffix(host, "/")
			}
		}
	}

	raw = strings.TrimPrefix(raw, "http://")
	raw = strings.TrimPrefix(raw, "https://")
	raw = strings.TrimPrefix(raw, "grpc://")
	raw = strings.TrimSuffix(raw, "/")
	return raw
}

func parseHeadersEnv(raw string) map[string]string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}

	headers := make(map[string]string)
	pairs := strings.Split(raw, ",")
	for _, pair := range pairs {
		trimmed := strings.TrimSpace(pair)
		if trimmed == "" {
			continue
		}
		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "" || value == "" {
			continue
		}
		headers[key] = value
	}
	if len(headers) == 0 {
		return nil
	}
	return headers
}
