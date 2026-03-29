package i18n

import (
	"strings"

	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Middleware injects the desired language into the request context before handling.
func Middleware(r *ghttp.Request) {
	lang := resolveLanguage(r)
	ctx := gi18n.WithLanguage(r.GetCtx(), lang)
	r.SetCtx(ctx)
	r.Middleware.Next()
}

func resolveLanguage(r *ghttp.Request) string {
	if lang := strings.TrimSpace(r.GetQuery("lang").String()); lang != "" {
		return normalizeLanguage(lang)
	}
	if lang := strings.TrimSpace(r.Header.Get("X-Language")); lang != "" {
		return normalizeLanguage(lang)
	}
	if accept := strings.TrimSpace(r.Header.Get("Accept-Language")); accept != "" {
		return normalizeLanguage(parseAcceptLanguage(accept))
	}
	return normalizeLanguage("")
}

// parseAcceptLanguage returns the first language token from the header.
func parseAcceptLanguage(header string) string {
	parts := strings.Split(header, ",")
	if len(parts) == 0 {
		return ""
	}
	primary := strings.SplitN(parts[0], ";", 2)[0]
	return strings.TrimSpace(primary)
}

func normalizeLanguage(lang string) string {
	lang = strings.TrimSpace(strings.ToLower(lang))
	switch lang {
	case "zh", "zh-cn":
		return "zh-CN"
	case "en", "en-us":
		return "en"
	default:
		return "zh-CN"
	}
}
