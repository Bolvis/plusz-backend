package env

import (
	"os"
	"testing"
)

func TestLoadDev(t *testing.T) {
	os.Setenv("ENV", "dev")

	Load()
	sslMode := os.Getenv("SSL_MODE")
	if sslMode != "disable" {
		t.Errorf(`SSL Mode should be set to "disable", acutall value: %v`, sslMode)
	}
}

func TestLoadProd(t *testing.T) {
	os.Setenv("ENV", "prod")

	Load()
	sslMode := os.Getenv("SSL_MODE")
	if sslMode != "require" {
		t.Errorf(`SSL Mode should be set to "require", acutall value: %v`, sslMode)
	}
}

func TestLoadNoValue(t *testing.T) {
	os.Setenv("ENV", "")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	Load()
}
