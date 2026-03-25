package response

import (
	"mime"
	"net/http"
	"strings"

	"github.com/TsingpekTao/shopa/iam-svc/internal/errs"
	"github.com/gogf/gf/v2/net/ghttp"
)

// HandlerResponse 閺勵垳绮烘稉鈧?HTTP 閸濆秴绨叉担鎾剁波閺嬪嫨鈧
type HandlerResponse struct {
	Code    int    `json:"code" dc:"娑撴艾濮熼柨娆掝嚖閻?`
	Message string `json:"message" dc:"閹绘劗銇氭穱鈩冧紖"`
	Data    any    `json:"data" dc:"閸濆秴绨查弫鐗堝祦"`
}

const (
	contentTypeEventStream  = "text/event-stream"
	contentTypeOctetStream  = "application/octet-stream"
	contentTypeMixedReplace = "multipart/x-mixed-replace"
)

var streamContentTypes = []string{contentTypeEventStream, contentTypeOctetStream, contentTypeMixedReplace}

// Middleware 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func Middleware(r *ghttp.Request) {
	if isDocRoute(r.URL.Path) {
		r.Middleware.Next()
		return
	}

	r.Middleware.Next()

	if r.Response.BufferLength() > 0 || r.Response.BytesWritten() > 0 {
		return
	}

	mediaType, _, _ := mime.ParseMediaType(r.Response.Header().Get("Content-Type"))
	for _, ct := range streamContentTypes {
		if mediaType == ct {
			return
		}
	}

	var (
		err = r.GetError()
		res = r.GetHandlerResponse()
	)
	if err == nil && r.Response.Status > 0 && r.Response.Status != http.StatusOK {
		switch r.Response.Status {
		case http.StatusForbidden, http.StatusUnauthorized:
			err = errs.New(errs.CodeAccessTokenInvalid)
		case http.StatusNotFound:
			err = errs.New(errs.CodeInvalidParam)
		default:
			err = errs.New(errs.CodeInternalError)
		}
		r.SetError(err)
	}

	code, message := errs.ToBiz(err)
	if err == nil {
		code = errs.CodeOK.Code()
		message = errs.CodeOK.Message()
	}

	r.Response.WriteJson(HandlerResponse{
		Code:    code,
		Message: message,
		Data:    res,
	})
}

// isDocRoute 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func isDocRoute(path string) bool {
	if path == "/api.json" {
		return true
	}
	if path == "/swagger" || path == "/swagger/" {
		return true
	}
	return strings.HasPrefix(path, "/swagger/")
}
