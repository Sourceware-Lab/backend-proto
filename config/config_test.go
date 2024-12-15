package config

import (
	"os"
	"testing"
)

func TestDbDSN_ParseDSN(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  string
		output DBDSN
	}{
		{
			name:  "valid DSN",
			input: "host=localhost port=5432 user=testUsername password=testPassword dbname=testdb sslmode=disable TimeZone=UTC",
			output: DBDSN{
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
			output: DBDSN{
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
			output: DBDSN{
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
			var d DBDSN
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
		input  DBDSN
		output string
	}{
		{
			name: "basic DSN",
			input: DBDSN{
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
			input:  DBDSN{},
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
	t.Parallel()
	tmpDir := t.TempDir()
	os.Setenv(ProjectName+"_PORT", "9090")
	os.Setenv(ProjectName+"_LOG_LEVEL", "info")
	os.Setenv(ProjectName+"_DATABASE_DSN", "host=localhost port=5432 user=test dbname=testdb sslmode=disable")
	os.Setenv(ProjectName+"_RELEASE_MODE", "true")
	os.Setenv(ProjectName+"_PROJECT_DIR", tmpDir)

	LoadConfig()

	if Config.Port != 9090 {
		t.Errorf("expected Port=9090, got %d", Config.Port)
	}
	if Config.LogLevel != "info" {
		t.Errorf("expected LogLevel=info, got %s", Config.LogLevel)
	}
	if Config.DatabaseDSN != "host=localhost port=5432 user=test dbname=testdb sslmode=disable" {
		t.Errorf("expected DatabaseDSN not matching, got %s", Config.DatabaseDSN)
	}
	if Config.ReleaseMode != true {
		t.Errorf("expected ReleaseMode=true, got %v", Config.ReleaseMode)
	}
	if Config.ProjectDir != tmpDir {
		t.Errorf("expected ProjectDir=%s, got %s", tmpDir, Config.ProjectDir)
	}
}
