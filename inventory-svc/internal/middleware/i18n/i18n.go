package i18n

import (
	"strings"

	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Middleware injects the normalized language into the request context for downstream handlers.
func Middleware(r *ghttp.Request) {
	lang := resolveLanguage(r)
	ctx := gi18n.WithLanguage(r.GetCtx(), lang)
	r.SetCtx(ctx)
	r.Middleware.Next()
}

func resolveLanguage(r *ghttp.Request) string {
	if langVar := r.GetQuery("lang"); langVar != nil {
		if lang := strings.TrimSpace(langVar.String()); lang != "" {
			return normalizeLanguage(lang)
		}
	}
	if lang := strings.TrimSpace(r.GetHeader("X-Language")); lang != "" {
		return normalizeLanguage(lang)
	}
	if header := r.GetHeader("Accept-Language"); header != "" {
		parts := strings.Split(header, ",")
		if len(parts) > 0 {
			langPart := strings.SplitN(strings.TrimSpace(parts[0]), ";", 2)[0]
			return normalizeLanguage(langPart)
		}
	}
	return normalizeLanguage("")
}

func normalizeLanguage(lang string) string {
	lang = strings.TrimSpace(lang)
	switch strings.ToLower(lang) {
	case "zh", "zh-cn":
		return "zh-CN"
	case "en", "en-us":
		return "en"
	default:
		return "zh-CN"
	}
}
