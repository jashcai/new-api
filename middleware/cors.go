package middleware

import (
	"net/url"
	"os"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return cors.New(BuildCORSConfig())
}

func BuildCORSConfig() cors.Config {
	config := cors.DefaultConfig()
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowCredentials = true
	config.AllowHeaders = allowedCORSHeaders()

	origins := allowedCORSOrigins()
	if len(origins) == 0 && !common.PrivateDeploymentMode {
		config.AllowAllOrigins = true
		config.AllowCredentials = false
		config.AllowHeaders = []string{"*"}
		return config
	}
	config.AllowOrigins = origins
	return config
}

func allowedCORSOrigins() []string {
	origins := make([]string, 0, len(common.CORSAllowedOrigins)+1)
	seen := map[string]struct{}{}
	add := func(origin string) {
		origin = normalizeOrigin(origin)
		if origin == "" {
			return
		}
		if _, ok := seen[origin]; ok {
			return
		}
		seen[origin] = struct{}{}
		origins = append(origins, origin)
	}
	for _, origin := range common.CORSAllowedOrigins {
		add(origin)
	}
	add(os.Getenv("FRONTEND_BASE_URL"))
	return origins
}

func normalizeOrigin(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "*" {
		return ""
	}
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return ""
	}
	parsed.Path = ""
	parsed.RawPath = ""
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return strings.TrimSuffix(parsed.String(), "/")
}

func allowedCORSHeaders() []string {
	if len(common.CORSAllowedHeaders) > 0 {
		return common.CORSAllowedHeaders
	}
	return []string{
		"Origin",
		"Content-Type",
		"Accept",
		"Authorization",
		"New-Api-User",
		"X-Api-Key",
		"X-Goog-Api-Key",
		"Anthropic-Version",
		"Mj-Api-Secret",
		"OpenAI-Organization",
		"OpenAI-Beta",
		common.RequestIdKey,
	}
}

func PoweredBy() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-New-Api-Version", common.Version)
		c.Next()
	}
}
