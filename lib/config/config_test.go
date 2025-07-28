package config

import (
	"os"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("Expected default host to be 0.0.0.0, got %s", cfg.Server.Host)
	}
	
	if cfg.Server.Port != 8080 {
		t.Errorf("Expected default port to be 8080, got %d", cfg.Server.Port)
	}
	
	if cfg.Log.Level != "info" {
		t.Errorf("Expected default log level to be info, got %s", cfg.Log.Level)
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name:    "valid config",
			config:  DefaultConfig(),
			wantErr: false,
		},
		{
			name: "invalid port",
			config: &Config{
				Server: ServerConfig{Port: -1},
				Raft:   RaftConfig{ElectionTimeout: time.Second, HeartbeatTimeout: time.Millisecond * 500},
				Log:    LogConfig{Level: "info", Format: "json"},
			},
			wantErr: true,
		},
		{
			name: "invalid log level",
			config: &Config{
				Server: ServerConfig{Port: 8080, GrpcTimeout: time.Second, ShutdownTimeout: time.Second},
				Raft:   RaftConfig{ElectionTimeout: time.Second, HeartbeatTimeout: time.Millisecond * 500},
				Log:    LogConfig{Level: "invalid", Format: "json"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadFromEnv(t *testing.T) {
	// Save original values
	originalHost := os.Getenv("NANOREDIZ_HOST")
	originalPort := os.Getenv("NANOREDIZ_PORT")
	originalLogLevel := os.Getenv("NANOREDIZ_LOG_LEVEL")
	
	// Clean up after test
	defer func() {
		os.Setenv("NANOREDIZ_HOST", originalHost)
		os.Setenv("NANOREDIZ_PORT", originalPort)
		os.Setenv("NANOREDIZ_LOG_LEVEL", originalLogLevel)
	}()
	
	// Set test environment variables
	os.Setenv("NANOREDIZ_HOST", "test-host")
	os.Setenv("NANOREDIZ_PORT", "9999")
	os.Setenv("NANOREDIZ_LOG_LEVEL", "debug")
	
	cfg := DefaultConfig()
	cfg.LoadFromEnv()
	
	if cfg.Server.Host != "test-host" {
		t.Errorf("Expected host to be test-host, got %s", cfg.Server.Host)
	}
	
	if cfg.Server.Port != 9999 {
		t.Errorf("Expected port to be 9999, got %d", cfg.Server.Port)
	}
	
	if cfg.Log.Level != "debug" {
		t.Errorf("Expected log level to be debug, got %s", cfg.Log.Level)
	}
}