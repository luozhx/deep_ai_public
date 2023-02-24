package main

import (
	"deep-ai-server/app/config"
	"deep-ai-server/app/db"
	authMiddleware "deep-ai-server/app/middleware/auth"
	"deep-ai-server/app/service/ai_model"
	"deep-ai-server/app/service/auth"
	"deep-ai-server/app/service/code"
	"deep-ai-server/app/service/dataset"
	"deep-ai-server/app/service/inference"
	"deep-ai-server/app/service/job"
	"deep-ai-server/app/service/system"
	"deep-ai-server/app/service/user"
	"deep-ai-server/app/tools/client"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db.Init(config.Database.Username, config.Database.Password,
		config.Database.Host, config.Database.Port, config.Database.Name)
	//client.Init("conf/kubeconfig_dev")
	client.Init("conf/kubeconfig")
	//db.DB().AutoMigrate(model.User{}, model.Job{}, model.Inference{}, model.Dataset{}, model.Code{}, model.AIModel{})

	e := echo.New()
	e.Debug = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(session.Middleware(sessions.NewCookieStore([]byte(config.App.SecretKey))))
	e.Use(authMiddleware.Middleware())
	//ce, _ := casbin.NewEnforcer("conf/casbin_model.conf", "conf/casbin_policy.csv")
	//e.Use(casbin_mw.Middleware(ce))

	r := e.Group("/api/v1")
	code.RegisterRouter(r)
	user.RegisterRouter(r)
	auth.RegisterRouter(r)
	dataset.RegisterRouter(r)
	inference.RegisterRouter(r)
	job.RegisterRouter(r)
	ai_model.RegisterRouter(r)
	system.RegisterRouter(r)

	e.Logger.Fatal(e.Start(":1323"))
}
