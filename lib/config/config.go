package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the nanorediz application
type Config struct {
	Server ServerConfig `json:"server"`
	Raft   RaftConfig   `json:"raft"`
	Log    LogConfig    `json:"log"`
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Host             string        `json:"host"`
	Port             int           `json:"port"`
	GrpcTimeout      time.Duration `json:"grpc_timeout"`
	ShutdownTimeout  time.Duration `json:"shutdown_timeout"`
	MaxConcurrentRPC int           `json:"max_concurrent_rpc"`
}

// RaftConfig holds Raft consensus algorithm configuration
type RaftConfig struct {
	ElectionTimeout  time.Duration `json:"election_timeout"`
	HeartbeatTimeout time.Duration `json:"heartbeat_timeout"`
	ApplyTimeout     time.Duration `json:"apply_timeout"`
	SnapshotRetain   int           `json:"snapshot_retain"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
	Output string `json:"output"`
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host:             "0.0.0.0",
			Port:             8080,
			GrpcTimeout:      30 * time.Second,
			ShutdownTimeout:  10 * time.Second,
			MaxConcurrentRPC: 100,
		},
		Raft: RaftConfig{
			ElectionTimeout:  1000 * time.Millisecond,
			HeartbeatTimeout: 500 * time.Millisecond,
			ApplyTimeout:     10 * time.Second,
			SnapshotRetain:   2,
		},
		Log: LogConfig{
			Level:  "info",
			Format: "json",
			Output: "stdout",
		},
	}
}

// LoadFromEnv loads configuration from environment variables
// Falls back to defaults if environment variables are not set
func (c *Config) LoadFromEnv() {
	if host := os.Getenv("NANOREDIZ_HOST"); host != "" {
		c.Server.Host = host
	}
	
	if portStr := os.Getenv("NANOREDIZ_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			c.Server.Port = port
		}
	}
	
	if timeoutStr := os.Getenv("NANOREDIZ_GRPC_TIMEOUT"); timeoutStr != "" {
		if timeout, err := time.ParseDuration(timeoutStr); err == nil {
			c.Server.GrpcTimeout = timeout
		}
	}
	
	if shutdownStr := os.Getenv("NANOREDIZ_SHUTDOWN_TIMEOUT"); shutdownStr != "" {
		if timeout, err := time.ParseDuration(shutdownStr); err == nil {
			c.Server.ShutdownTimeout = timeout
		}
	}
	
	if logLevel := os.Getenv("NANOREDIZ_LOG_LEVEL"); logLevel != "" {
		c.Log.Level = logLevel
	}
	
	if logFormat := os.Getenv("NANOREDIZ_LOG_FORMAT"); logFormat != "" {
		c.Log.Format = logFormat
	}
	
	if logOutput := os.Getenv("NANOREDIZ_LOG_OUTPUT"); logOutput != "" {
		c.Log.Output = logOutput
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}
	
	if c.Server.GrpcTimeout <= 0 {
		return fmt.Errorf("grpc timeout must be positive")
	}
	
	if c.Server.ShutdownTimeout <= 0 {
		return fmt.Errorf("shutdown timeout must be positive")
	}
	
	if c.Raft.ElectionTimeout <= 0 {
		return fmt.Errorf("election timeout must be positive")
	}
	
	if c.Raft.HeartbeatTimeout <= 0 {
		return fmt.Errorf("heartbeat timeout must be positive")
	}
	
	validLogLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true,
	}
	if !validLogLevels[c.Log.Level] {
		return fmt.Errorf("invalid log level: %s", c.Log.Level)
	}
	
	validLogFormats := map[string]bool{
		"json": true, "text": true,
	}
	if !validLogFormats[c.Log.Format] {
		return fmt.Errorf("invalid log format: %s", c.Log.Format)
	}
	
	return nil
}