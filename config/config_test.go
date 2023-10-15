package config_test

import (
	"testing"

	"github.com/VassilisPallas/gvs/config"
)

func TestConfig(t *testing.T) {
	baseURL := "https://go.dev/dl"
	requestTimeout := 30
	expireCacheAfter := float64(24 * 7)

	cf := config.GetConfig()

	if cf.GO_BASE_URL != baseURL {
		t.Errorf("GO_BASE_URL should be %q, instead got %q", baseURL, cf.GO_BASE_URL)
	}

	if cf.REQUEST_TIMEOUT != requestTimeout {
		t.Errorf("REQUEST_TIMEOUT should be %q, instead got %q", requestTimeout, cf.REQUEST_TIMEOUT)
	}

	if cf.EXPIRE_CACHE_AFTER != expireCacheAfter {
		t.Errorf("EXPIRE_CACHE_AFTER should be %.1f, instead got %.1f", expireCacheAfter, cf.EXPIRE_CACHE_AFTER)
	}
}
