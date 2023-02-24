package user

import (
	"deep-ai-server/app/db"
	"deep-ai-server/app/model"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

func RegisterRouter(e *echo.Group) {
	r := e.Group("/user")

	r.GET("", getList)
	r.GET("/:id", get)
	r.POST("", post)
	r.PUT("", put)
	r.DELETE("/:id", _delete)
}

func getList(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, "")
}

func get(ctx echo.Context) error {
	userID := ctx.Get("id").(uint)
	u := new(model.User)
	db.DB().First(u, userID)
	return ctx.JSON(http.StatusOK, u)
}

func post(ctx echo.Context) error {
	data := struct {
		Username string
		Password string
		Email    string
	}{}
	_ = ctx.Bind(&data)
	pwhash, _ := bcrypt.GenerateFromPassword([]byte(data.Password), 14)
	u := model.User{Username: data.Username, PasswordHash: string(pwhash), Email: data.Email, Role: model.USER_ROLE_USER}
	db.DB().Create(&u)

	return ctx.JSON(http.StatusCreated, echo.Map{"id": u.ID})
}

func put(ctx echo.Context) error {
	u := ctx.Get("user").(*model.User)
	data := struct {
		Password string
	}{}
	_ = ctx.Bind(&data)
	pwhash, _ := bcrypt.GenerateFromPassword([]byte(data.Password), 14)
	db.DB().Model(u).Update(model.User{PasswordHash: string(pwhash)})
	return ctx.NoContent(http.StatusOK)
}

func _delete(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))
	db.DB().Delete(&model.User{ID: uint(id)})
	return ctx.NoContent(http.StatusNoContent)
}
