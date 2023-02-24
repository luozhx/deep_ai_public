package inference

import (
	"deep-ai-server/app/db"
	"deep-ai-server/app/model"
	"deep-ai-server/app/tools/client"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httputil"
	"strconv"
)

func RegisterRouter(e *echo.Group) {
	r := e.Group("/inference")

	r.GET("", getList)
	r.GET("/:id", get)
	r.POST("", post)
	r.PUT("/:id", put)
	r.DELETE("/:id", _delete)

	r.POST("/:id/inference/:service_name", proxy)
}

func getList(ctx echo.Context) error {
	u := ctx.Get("user").(*model.User)
	db.DB().Model(u).Where(&model.Inference{Framework: ctx.QueryParam("framework")}).Related(&u.Inferences)
	return ctx.JSON(http.StatusOK, u.Inferences)
}

func get(ctx echo.Context) error {
	u := ctx.Get("user").(*model.User)
	inferenceID, _ := strconv.Atoi(ctx.Param("id"))

	inference := new(model.Inference)
	db.DB().Where(&model.Inference{UserID: u.ID}).First(inference, inferenceID)

	// TODO: 获取service状态

	return ctx.JSON(http.StatusOK, inference)
}

func post(ctx echo.Context) error {
	u := ctx.Get("user").(*model.User)
	// 参数
	data := struct {
		Name             string
		Description      string
		Framework        string
		FrameworkVersion string
		//Num              int
		//CPU              int
		//Memory           int
		//GPU              int
		AIModelID uint
		Dimension string
		ModelFile string
	}{}
	_ = ctx.Bind(&data)

	AIModel := new(model.AIModel)
	if db.DB().Where(&model.AIModel{UserID: u.ID}).First(AIModel, data.AIModelID).RecordNotFound() {
		return echo.ErrBadRequest
	}

	inference := model.Inference{
		Name:             data.Name,
		Description:      data.Description,
		Framework:        data.Framework,
		FrameworkVersion: data.FrameworkVersion,
		Dimension:        data.Dimension,
		ModelFile:        data.ModelFile,
		//Num:              data.Num,
		//CPU:              data.CPU,
		//Memory:           data.Memory,
		//GPU:              data.GPU,
		UserID:    u.ID,
		AIModelID: data.AIModelID,
	}
	db.DB().Create(&inference)

	stopListener := make(chan struct{})
	go createServing(inference, *AIModel, stopListener)

	return ctx.JSON(http.StatusCreated, echo.Map{"id": inference.ID})
}

func put(ctx echo.Context) error {
	u := ctx.Get("user").(*model.User)
	inferenceID, _ := strconv.Atoi(ctx.Param("id"))
	data := struct {
		Name        string
		Description string
	}{}
	_ = ctx.Bind(&data)
	db.DB().Model(&model.Inference{ID: uint(inferenceID), UserID: u.ID}).Update(&model.Inference{Name: data.Name, Description: data.Description})

	return ctx.NoContent(http.StatusOK)
}

func _delete(ctx echo.Context) error {
	u := ctx.Get("user").(*model.User)
	inferenceID, _ := strconv.Atoi(ctx.Param("id"))

	cli := client.NewClient(client.NAMESPACE_DEEP_AI)
	_ = cli.DeleteServing(fmt.Sprintf("user%d-serving%d", u.ID, inferenceID))
	db.DB().Delete(&model.Inference{ID: uint(inferenceID), UserID: u.ID})
	return ctx.NoContent(http.StatusNoContent)
}

// reference: https://www.integralist.co.uk/posts/golang-reverse-proxy/
func proxy(c echo.Context) error {
	//todo: verify user
	name := c.Param("service_name")
	handler := httputil.ReverseProxy{}
	handler.Director = func(request *http.Request) {
		request.URL.Scheme = "http"
		request.URL.Host = fmt.Sprintf("%s.deep-ai.svc.cluster.local", name)
		request.URL.Path = "/v1/models/model:predict"
		if _, ok := request.Header["User-Agent"]; !ok {
			request.Header.Set("User-Agent", "")
		}
	}
	handler.ServeHTTP(c.Response(), c.Request())
	return nil
}
