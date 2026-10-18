package config

import "testing"

func TestLoadOTELDefaults(t *testing.T) {
	t.Setenv("SCANNER_URL", "http://scanner:8081")
	t.Setenv("API_KEY", "test-key")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.ServiceName != "api-gateway" {
		t.Errorf("ServiceName = %q, want %q", cfg.ServiceName, "api-gateway")
	}
	if cfg.ServiceVersion != "0.1.0" {
		t.Errorf("ServiceVersion = %q, want %q", cfg.ServiceVersion, "0.1.0")
	}
	if cfg.OTELExporterOTLPEndpoint != "" {
		t.Errorf("OTELExporterOTLPEndpoint = %q, want empty", cfg.OTELExporterOTLPEndpoint)
	}
}

func TestLoadOTELConfiguration(t *testing.T) {
	t.Setenv("SCANNER_URL", "http://scanner:8081")
	t.Setenv("API_KEY", "test-key")
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "https://otel.example.com:443")
	t.Setenv("OTEL_EXPORTER_OTLP_HEADERS", "Authorization=Basic token")
	t.Setenv("SERVICE_NAME", "gateway-test")
	t.Setenv("SERVICE_VERSION", "1.2.3")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.OTELExporterOTLPEndpoint != "https://otel.example.com:443" {
		t.Errorf("OTELExporterOTLPEndpoint = %q", cfg.OTELExporterOTLPEndpoint)
	}
	if cfg.OTELExporterOTLPHeaders != "Authorization=Basic token" {
		t.Errorf("OTELExporterOTLPHeaders = %q", cfg.OTELExporterOTLPHeaders)
	}
	if cfg.ServiceName != "gateway-test" {
		t.Errorf("ServiceName = %q, want %q", cfg.ServiceName, "gateway-test")
	}
	if cfg.ServiceVersion != "1.2.3" {
		t.Errorf("ServiceVersion = %q, want %q", cfg.ServiceVersion, "1.2.3")
	}
}
