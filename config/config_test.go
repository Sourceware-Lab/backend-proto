package config_test

import (
	"testing"

	"github.com/Sourceware-Lab/go-huma-gin-postgres-template/config"
)

func TestDbDSN_ParseDSN(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  string
		output config.DBDSN
	}{
		{
			name:  "valid DSN",
			input: "host=localhost port=5432 user=testUsername password=testPassword dbname=testdb sslmode=disable TimeZone=UTC",
			output: config.DBDSN{
				Host:     "localhost",
				Port:     5432,
				User:     "testUsername",
				Password: "testPassword",
				DBName:   "testdb",
				SSLMode:  "disable",
				TimeZone: "UTC",
			},
		},
		{
			name:  "missing optional fields in DSN",
			input: "host=127.0.0.1 port=3306 user=root dbname=appdb sslmode=required",
			output: config.DBDSN{
				Host:     "127.0.0.1",
				Port:     3306,
				User:     "root",
				Password: "",
				DBName:   "appdb",
				SSLMode:  "required",
				TimeZone: "",
			},
		},
		{
			name:  "empty DSN",
			input: "",
			output: config.DBDSN{
				Host:     "",
				Port:     0,
				User:     "",
				Password: "",
				DBName:   "",
				SSLMode:  "",
				TimeZone: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var d config.DBDSN
			result := d.ParseDSN(tt.input)

			if result != tt.output {
				t.Errorf("ParseDSN(%q) = %v, want %v", tt.input, result, tt.output)
			}
		})
	}
}

func TestDbDSN_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  config.DBDSN
		output string
	}{
		{
			name: "basic DSN",
			input: config.DBDSN{
				Host:     "localhost",
				Port:     5432,
				User:     "testUsername",
				Password: "testPassword",
				DBName:   "testdb",
				SSLMode:  "disable",
				TimeZone: "UTC",
			},
			output: "host=localhost user=testUsername password=testPassword dbname=testdb " +
				"port=5432 sslmode=disable TimeZone=UTC",
		},
		{
			name:   "empty DSN",
			input:  config.DBDSN{},
			output: "host= user= password= dbname= port=0 sslmode= TimeZone=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := tt.input.String()
			if result != tt.output {
				t.Errorf("String() = %q, want %q", result, tt.output)
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("PORT", "9090")
	t.Setenv("LOG_LEVEL", "info")
	t.Setenv("DATABASE_DSN", "host=localhost port=5432 user=test dbname=testdb sslmode=disable")
	t.Setenv("RELEASE_MODE", "true")
	t.Setenv("PROJECT_DIR", tmpDir)

	config.LoadConfig()

	if config.Config.Port != 9090 {
		t.Errorf("expected Port=9090, got %d", config.Config.Port)
	}

	if config.Config.LogLevel != "info" {
		t.Errorf("expected LogLevel=info, got %s", config.Config.LogLevel)
	}

	if config.Config.DatabaseDSN != "host=localhost port=5432 user=test dbname=testdb sslmode=disable" {
		t.Errorf("expected DatabaseDSN not matching, got %s", config.Config.DatabaseDSN)
	}

	if config.Config.ReleaseMode != true {
		t.Errorf("expected ReleaseMode=true, got %v", config.Config.ReleaseMode)
	}

	if config.Config.ProjectDir != tmpDir {
		t.Errorf("expected ProjectDir=%s, got %s", tmpDir, config.Config.ProjectDir)
	}
}
