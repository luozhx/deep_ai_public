package system

import (
	"deep-ai-server/app/db"
	"deep-ai-server/app/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

func RegisterRouter(e *echo.Group) {
	r := e.Group("/system")

	r.GET("/status", get)
}

func get(ctx echo.Context) error {
	response := echo.Map{}
	type Status struct {
		Status int `json:"status"`
		Count  int `json:"count"`
	}
	items := map[string]interface{}{
		"code":      &model.Code{},
		"dataset":   &model.Dataset{},
		"model":     &model.AIModel{},
		"job":       &model.Job{},
		"inference": &model.Inference{},
	}
	for k, v := range items {
		status := make([]Status, 0)
		db.DB().Model(v).Select("status, count(*) as count").Group("status").Scan(&status)
		response[k] = status
	}
	
	return ctx.JSON(http.StatusOK, response)
}
