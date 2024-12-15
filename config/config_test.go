package config

import (
	"fmt"
	"os"
	"testing"
)

func TestDbDSN_ParseDSN(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output DbDSN
	}{
		{
			name:  "valid DSN",
			input: "host=localhost port=5432 user=testUsername password=testPassword dbname=testdb sslmode=disable TimeZone=UTC",
			output: DbDSN{
				Host:     "localhost",
				Port:     5432,
				User:     "testUsername",
				Password: "testPassword",
				DbName:   "testdb",
				SSLMode:  "disable",
				TimeZone: "UTC",
			},
		},
		{
			name:  "missing optional fields in DSN",
			input: "host=127.0.0.1 port=3306a user=root dbname=appdb sslmode=required",
			output: DbDSN{
				Host:     "127.0.0.1",
				Port:     3306,
				User:     "root",
				Password: "",
				DbName:   "appdb",
				SSLMode:  "required",
				TimeZone: "",
			},
		},
		{
			name:  "empty DSN",
			input: "",
			output: DbDSN{
				Host:     "",
				Port:     0,
				User:     "",
				Password: "",
				DbName:   "",
				SSLMode:  "",
				TimeZone: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d DbDSN
			result := d.ParseDSN(tt.input)
			if result != tt.output {
				t.Errorf("ParseDSN(%q) = %v, want %v", tt.input, result, tt.output)
			}
		})
	}
}

func TestDbDSN_String(t *testing.T) {
	tests := []struct {
		name   string
		input  DbDSN
		output string
	}{
		{
			name: "basic DSN",
			input: DbDSN{
				Host:     "localhost",
				Port:     5432,
				User:     "testUsername",
				Password: "testPassword",
				DbName:   "testdb",
				SSLMode:  "disable",
				TimeZone: "UTC",
			},
			output: "host=localhost user=testUsername password=testPassword dbname=testdb port=5432 sslmode=disable TimeZone=UTC",
		},
		{
			name:   "empty DSN",
			input:  DbDSN{},
			output: "host= user= password= dbname= port=0 sslmode= TimeZone=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.String()
			if result != tt.output {
				t.Errorf("String() = %q, want %q", result, tt.output)
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv(fmt.Sprintf("%s_PORT", ProjectName), "9090")
	os.Setenv(fmt.Sprintf("%s_LOG_LEVEL", ProjectName), "info")
	os.Setenv(fmt.Sprintf("%s_DATABASE_DSN", ProjectName), "host=localhost port=5432 user=test dbname=testdb sslmode=disable")
	os.Setenv(fmt.Sprintf("%s_RELEASE_MODE", ProjectName), "true")
	os.Setenv(fmt.Sprintf("%s_PROJECT_DIR", ProjectName), tmpDir)

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
