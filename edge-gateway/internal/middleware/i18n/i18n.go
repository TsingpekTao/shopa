package i18n

import (
	"strings"

	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/net/ghttp"
)

const defaultLanguage = "zh-CN"

// Middleware injects request language into context for i18n translation.
func Middleware(r *ghttp.Request) {
	var language = strings.TrimSpace(r.Get("lang").String())
	if language == "" {
		language = strings.TrimSpace(r.GetHeader("X-Language"))
	}
	if language == "" {
		language = firstAcceptLanguage(strings.TrimSpace(r.GetHeader("Accept-Language")))
	}
	language = normalizeLanguage(language)
	r.SetCtx(gi18n.WithLanguage(r.GetCtx(), language))
	r.Middleware.Next()
}

func firstAcceptLanguage(raw string) string {
	if raw == "" {
		return ""
	}
	var first = strings.TrimSpace(strings.Split(raw, ",")[0])
	if index := strings.Index(first, ";"); index >= 0 {
		first = strings.TrimSpace(first[:index])
	}
	return first
}

func normalizeLanguage(raw string) string {
	var language = strings.ToLower(strings.TrimSpace(raw))
	switch {
	case language == "":
		return defaultLanguage
	case strings.HasPrefix(language, "zh"):
		return "zh-CN"
	case strings.HasPrefix(language, "en"):
		return "en"
	default:
		return defaultLanguage
	}
}
