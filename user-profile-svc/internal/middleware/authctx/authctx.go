package authctx

import (
	"net/http"

	"github.com/TsingpekTao/shopa/user-profile-svc/internal/ctxkey"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
)

// UserIDFromHeader 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func UserIDFromHeader(r *ghttp.Request) {
	// Read user id from common headers.
	// In real production, this should be replaced by JWT verification (iam-svc) + user id claim.
	var (
		raw = r.GetHeader("X-User-Id")
	)
	if raw == "" {
		raw = r.GetHeader("User-Id")
	}
	if raw == "" {
		raw = r.GetHeader("Uid")
	}
	if raw == "" {
		raw = r.GetHeader("X-Uid")
	}

	// Reject requests without user identity to prevent accidental anonymous access to "me/*" APIs.
	if raw == "" {
		// Set HTTP status for client side.
		r.Response.WriteStatus(http.StatusUnauthorized)
		// Attach error for unified response middleware + error logging.
		r.SetError(gerror.NewCode(gcode.CodeNotAuthorized, "missing X-User-Id header"))
		// Stop the request pipeline.
		r.Exit()
	}

	// Convert to uint64 and validate.
	userID := gconv.Uint64(raw)
	if userID == 0 {
		// Use gerror so it is logged with stack in GoFrame error logger.
		r.SetError(gerror.NewCode(gcode.CodeInvalidParameter, "invalid X-User-Id header"))
		r.Response.WriteStatus(http.StatusBadRequest)
		// Stop the request pipeline.
		r.Exit()
	}

	// Store user id into request context for downstream handlers.
	r.SetCtxVar(ctxkey.UserIDKey{}, userID)

	// Continue to next middleware/handler.
	r.Middleware.Next()
}
