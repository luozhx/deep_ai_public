package auth

import (
	"deep-ai-server/app/model"
	"github.com/labstack/echo/v4"
)

func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			//sess, _ := session.Get("session", ctx)
			//id, username, role := uint(0), "", 0
			//if sess.Values["id"] != nil && sess.Values["username"] != nil && sess.Values["role"] != nil {
			//	id = sess.Values["id"].(uint)
			//	username = sess.Values["username"].(string)
			//	role = sess.Values["role"].(int)
			//}
			//ctx.Set("user", &model.User{
			//	ID:                        id,
			//	Username:                  username,
			//	Role:                      role,
			//})
			ctx.Set("user", &model.User{
				ID:                        1,
				Username:                  "admin",
				Role:                      0,
			})
			return next(ctx)
		}
	}
}
