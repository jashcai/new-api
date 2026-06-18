package middleware

import (
	"reflect"
	"testing"

	"github.com/QuantumNous/new-api/common"
)

func TestBuildCORSConfigDefaultAllowsAllOrigins(t *testing.T) {
	oldPrivateMode := common.PrivateDeploymentMode
	oldOrigins := common.CORSAllowedOrigins
	oldHeaders := common.CORSAllowedHeaders
	t.Cleanup(func() {
		common.PrivateDeploymentMode = oldPrivateMode
		common.CORSAllowedOrigins = oldOrigins
		common.CORSAllowedHeaders = oldHeaders
	})

	common.PrivateDeploymentMode = false
	common.CORSAllowedOrigins = nil
	common.CORSAllowedHeaders = nil

	config := BuildCORSConfig()
	if !config.AllowAllOrigins {
		t.Fatal("expected default CORS config to allow all origins for compatibility")
	}
	if config.AllowCredentials {
		t.Fatal("expected wildcard CORS config to disable credentials")
	}
	if !reflect.DeepEqual(config.AllowHeaders, []string{"*"}) {
		t.Fatalf("expected wildcard headers, got %#v", config.AllowHeaders)
	}
}

func TestBuildCORSConfigPrivateModeRestrictsOrigins(t *testing.T) {
	oldPrivateMode := common.PrivateDeploymentMode
	oldOrigins := common.CORSAllowedOrigins
	oldHeaders := common.CORSAllowedHeaders
	t.Cleanup(func() {
		common.PrivateDeploymentMode = oldPrivateMode
		common.CORSAllowedOrigins = oldOrigins
		common.CORSAllowedHeaders = oldHeaders
	})
	t.Setenv("FRONTEND_BASE_URL", "https://console.example.com/app")

	common.PrivateDeploymentMode = true
	common.CORSAllowedOrigins = []string{
		"https://gateway.example.com",
		"https://gateway.example.com/path?x=1",
		"*",
		"not-a-url",
	}
	common.CORSAllowedHeaders = []string{"Authorization", "Content-Type"}

	config := BuildCORSConfig()
	if config.AllowAllOrigins {
		t.Fatal("expected private CORS config to disable allow-all origins")
	}
	expectedOrigins := []string{
		"https://gateway.example.com",
		"https://console.example.com",
	}
	if !reflect.DeepEqual(config.AllowOrigins, expectedOrigins) {
		t.Fatalf("expected origins %#v, got %#v", expectedOrigins, config.AllowOrigins)
	}
	if !reflect.DeepEqual(config.AllowHeaders, common.CORSAllowedHeaders) {
		t.Fatalf("expected configured headers, got %#v", config.AllowHeaders)
	}
}

func TestBuildCORSConfigPrivateModeWithNoOriginsAllowsOnlySameOrigin(t *testing.T) {
	oldPrivateMode := common.PrivateDeploymentMode
	oldOrigins := common.CORSAllowedOrigins
	t.Cleanup(func() {
		common.PrivateDeploymentMode = oldPrivateMode
		common.CORSAllowedOrigins = oldOrigins
	})
	t.Setenv("FRONTEND_BASE_URL", "")

	common.PrivateDeploymentMode = true
	common.CORSAllowedOrigins = nil

	config := BuildCORSConfig()
	if config.AllowAllOrigins {
		t.Fatal("expected private CORS config to disable allow-all origins")
	}
	if len(config.AllowOrigins) != 0 {
		t.Fatalf("expected no cross-origin allowlist, got %#v", config.AllowOrigins)
	}
}
