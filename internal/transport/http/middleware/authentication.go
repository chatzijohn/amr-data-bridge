package middleware

import (
	"amr-data-bridge/config"
	"log/slog"
	"net"
	"net/http"
	"strings"
)

func Auth(cfg config.AuthConfig, tokens config.AuthTokens, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerName := cfg.Header
		if headerName == "" {
			headerName = "Authorization"
		}

		val := r.Header.Get(headerName)
		if val == "" {
			slog.Warn("Auth failed: missing header", "header", headerName, "remote_addr", r.RemoteAddr)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var token string
		// If using standard Authorization header, expect Bearer prefix.
		// Custom headers (like X-API-Key) usually contain just the token.
		switch headerName {
		case "Authorization":
			const prefix = "Bearer "
			if len(val) < len(prefix) || val[:len(prefix)] != prefix {
				slog.Warn("Auth failed: invalid bearer prefix", "remote_addr", r.RemoteAddr)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			token = val[len(prefix):]
		default:
			token = val
		}

		// 1. Identify the token name (key) by matching the provided token value against loaded secrets
		var tokenName string
		for name, secret := range tokens {
			if secret == token {
				tokenName = name
				break
			}
		}

		if tokenName == "" {
			slog.Warn("Auth failed: token not found", "remote_addr", r.RemoteAddr)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// 2. Check if this token name has a policy in preferences.yaml
		policy, exists := cfg.Tokens[tokenName]
		if !exists {
			slog.Warn("Auth failed: no policy defined for token", "token_name", tokenName, "remote_addr", r.RemoteAddr)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		if !policy.Enabled {
			slog.Warn("Auth failed: token disabled", "token_name", tokenName)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Check IP Whitelist if defined
		if len(policy.IPs) > 0 {
			// Determine Client IP (handling proxies)
			var host string

			if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
				// X-Forwarded-For can contain multiple IPs, the first is the client
				parts := strings.Split(xff, ",")
				host = strings.TrimSpace(parts[0])
			} else if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
				host = xrip
			} else {
				var err error
				host, _, err = net.SplitHostPort(r.RemoteAddr)
				if err != nil {
					host = r.RemoteAddr
				}
			}

			allowed := false
			clientIP := net.ParseIP(host)

			for _, cidr := range policy.IPs {
				// 1. Try parsing as CIDR (e.g. 192.168.1.0/24)
				_, ipNet, err := net.ParseCIDR(cidr)
				if err == nil && clientIP != nil && ipNet.Contains(clientIP) {
					allowed = true
					break
				}
				// 2. Fallback to exact IP match
				if cidr == host {
					allowed = true
					break
				}
			}

			if !allowed {
				slog.Warn("Auth failed: IP not whitelisted", "token_name", tokenName, "client_ip", host, "allowed_ips", policy.IPs)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
