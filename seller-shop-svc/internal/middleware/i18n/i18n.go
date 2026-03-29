package i18n

import (
	"strings"

	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/net/ghttp"
)

func Middleware(r *ghttp.Request) {
	lang := firstNonEmpty(
		r.GetQuery("lang").String(),
		r.GetHeader("X-Language"),
		acceptLanguageLang(r.GetHeader("Accept-Language")),
	)

	lang = normalizeLanguage(lang)
	r.SetCtx(gi18n.WithLanguage(r.GetCtx(), lang))

	r.Middleware.Next()
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if trimmed := strings.TrimSpace(v); trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func acceptLanguageLang(header string) string {
	if trimmed := strings.TrimSpace(header); trimmed != "" {
		if idx := strings.Index(trimmed, ","); idx >= 0 {
			trimmed = trimmed[:idx]
		}
		if idx := strings.Index(trimmed, ";"); idx >= 0 {
			trimmed = trimmed[:idx]
		}
		return strings.TrimSpace(trimmed)
	}
	return ""
}

func normalizeLanguage(lang string) string {
	switch strings.ToLower(strings.TrimSpace(lang)) {
	case "zh", "zh-cn":
		return "zh-CN"
	case "en", "en-us":
		return "en"
	default:
		return "zh-CN"
	}
}
