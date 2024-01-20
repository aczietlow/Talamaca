package config

import (
	"github.com/spf13/afero"
	"reflect"
	"testing"
)

func TestLoadConfigValidFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	fileName := "/tmp/config.json"
	content := `{"token":"abc123","user":"testuser","owner":"testowner","repo":"testrepo","users":["user1", "user2"]}`

	afero.WriteFile(fs, fileName, []byte(content), 0644)

	oldOsReadFile := osReadFile
	osReadFile = func(name string) ([]byte, error) {
		return afero.ReadFile(fs, name)
	}
	defer func() { osReadFile = oldOsReadFile }()

	expectedConfig := &Config{
		Token:    "abc123",
		Username: "testuser",
		Owner:    "testowner",
		Repo:     "testrepo",
		Users:    []string{"user1", "user2"},
	}

	config, err := LoadConfig(fileName)
	if err != nil {
		t.Fatalf("Failed to load config: %s", err)
	}

	if !reflect.DeepEqual(config, expectedConfig) {
		t.Errorf("Expected %+v, got %+v", expectedConfig, config)
	}
}

func TestLoadConfigNonExistentFile(t *testing.T) {
	_, err := LoadConfig("/path/to/nonexistent/file")
	if err == nil {
		t.Errorf("Expected an error for non-existent file, got none")
	}
}

func TestLoadConfigInvalidJSON(t *testing.T) {
	fs := afero.NewMemMapFs()
	fileName := "/tmp/invalid.json"
	content := `{"token": "abc123"`

	afero.WriteFile(fs, fileName, []byte(content), 0644)

	oldOsReadFile := osReadFile
	osReadFile = func(name string) ([]byte, error) {
		return afero.ReadFile(fs, name)
	}
	defer func() { osReadFile = oldOsReadFile }()

	_, err := LoadConfig(fileName)
	if err == nil {
		t.Errorf("Expected JSON unmarshal error, got nil")
	}
}
