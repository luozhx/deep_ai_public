package auth

import (
	"deep-ai-server/app/db"
	"deep-ai-server/app/model"
	"deep-ai-server/app/service/user"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
)

func RegisterRouter(e *echo.Group) {
	r := e.Group("/auth")

	r.POST("", post)
	r.DELETE("", _delete)
}

func post(ctx echo.Context) error {
	data := struct {
		UsernameOrEmail string
		Password        string
	}{}
	_ = ctx.Bind(&data)
	u := new(model.User)
	if db.DB().Model(u).Where(model.User{Username: data.UsernameOrEmail}).Or(model.User{Email: data.UsernameOrEmail}).First(u).RecordNotFound() {
		return echo.ErrUnauthorized
	}
	if !user.IsPasswordValid(u, data.Password) {
		return echo.ErrUnauthorized
	}
	user.Login(u)

	sess, _ := session.Get("session", ctx)
	sess.Options = &sessions.Options{
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Values["id"] = u.ID
	sess.Values["username"] = u.Username
	sess.Values["role"] = u.Role
	_ = sess.Save(ctx.Request(), ctx.Response())

	return ctx.NoContent(http.StatusCreated)
}

func _delete(ctx echo.Context) error { //TESTED
	sess, _ := session.Get("session", ctx)
	delete(sess.Values, "id")
	delete(sess.Values, "username")
	delete(sess.Values, "role")
	delete(sess.Values, "pv")
	_ = sess.Save(ctx.Request(), ctx.Response())
	return ctx.NoContent(http.StatusNoContent)
}
