package i18n

import (
	"strings"

	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/net/ghttp"
)

const defaultLang = "zh-CN"

// Middleware injects the selected language into the request context so downstream
// handlers can rely on g.I18n().T().
func Middleware(r *ghttp.Request) {
	lang := detectLanguage(r)
	ctx := gi18n.WithLanguage(r.GetCtx(), lang)
	r.SetCtx(ctx)
	r.Middleware.Next()
}

func detectLanguage(r *ghttp.Request) string {
	for _, source := range []string{
		r.GetQuery("lang").String(),
		r.Header.Get("X-Language"),
		acceptLanguage(r.Header.Get("Accept-Language")),
	} {
		if normalized := normalizeLanguage(source); normalized != "" {
			return normalized
		}
	}
	return defaultLang
}

func acceptLanguage(header string) string {
	if header == "" {
		return ""
	}
	parts := strings.Split(header, ",")
	if len(parts) == 0 {
		return ""
	}
	first := strings.TrimSpace(parts[0])
	if idx := strings.Index(first, ";"); idx > -1 {
		first = first[:idx]
	}
	return first
}

func normalizeLanguage(lang string) string {
	switch strings.ToLower(strings.TrimSpace(lang)) {
	case "zh", "zh-cn":
		return "zh-CN"
	case "en", "en-us":
		return "en"
	default:
		return ""
	}
}
