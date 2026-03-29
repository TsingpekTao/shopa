package i18n

import (
	"strings"

	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/text/gstr"
)

const defaultLanguage = "zh-CN"

var languageMap = map[string]string{
	"zh":    "zh-CN",
	"zh-cn": "zh-CN",
	"en":    "en",
	"en-us": "en",
}

// Middleware sets the current language in the request context and continues the chain.
func Middleware(r *ghttp.Request) {
	lang := detectLanguage(r)
	r.SetCtx(gi18n.WithLanguage(r.GetCtx(), lang))
	r.Middleware.Next()
}

func detectLanguage(r *ghttp.Request) string {
	queryLang := ""
	if q := r.GetQuery("lang"); q != nil {
		queryLang = q.String()
	}

	sources := []string{
		queryLang,
		r.Header.Get("X-Language"),
		parseAcceptLanguage(r.Header.Get("Accept-Language")),
	}
	for _, candidate := range sources {
		if candidate = strings.TrimSpace(candidate); candidate != "" {
			return normalizeLanguage(candidate)
		}
	}
	return defaultLanguage
}

func normalizeLanguage(raw string) string {
	raw = strings.TrimSpace(raw)
	if idx := strings.Index(raw, ";"); idx != -1 {
		raw = raw[:idx]
	}
	raw = strings.ToLower(strings.TrimSpace(raw))
	if normalized, ok := languageMap[raw]; ok {
		return normalized
	}
	return defaultLanguage
}

func parseAcceptLanguage(header string) string {
	header = strings.TrimSpace(header)
	if header == "" {
		return ""
	}
	parts := gstr.SplitAndTrim(header, ",")
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}
