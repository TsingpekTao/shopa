package i18n

import (
	"strings"

	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Middleware selects the best supported language and injects it into the request context.
func Middleware(r *ghttp.Request) {
	lang := pickLanguage(r)
	r.SetCtx(gi18n.WithLanguage(r.GetCtx(), lang))
	r.Middleware.Next()
}

func pickLanguage(r *ghttp.Request) string {
	for _, extractor := range []func(*ghttp.Request) string{
		func(req *ghttp.Request) string { return req.GetQuery("lang").String() },
		func(req *ghttp.Request) string { return req.GetHeader("X-Language") },
		func(req *ghttp.Request) string { return firstAcceptLanguage(req.GetHeader("Accept-Language")) },
	} {
		if value := normalizeLanguage(extractor(r)); value != "" {
			return value
		}
	}
	return "zh-CN"
}

func firstAcceptLanguage(header string) string {
	if header == "" {
		return ""
	}
	for _, part := range strings.Split(header, ",") {
		lang := strings.TrimSpace(strings.Split(part, ";")[0])
		if lang != "" {
			return lang
		}
	}
	return ""
}

func normalizeLanguage(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "zh", "zh-cn":
		return "zh-CN"
	case "en", "en-us":
		return "en"
	default:
		return ""
	}
}
