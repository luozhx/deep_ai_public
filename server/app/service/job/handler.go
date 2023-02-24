package job

import (
	"deep-ai-server/app/db"
	"deep-ai-server/app/model"
	"deep-ai-server/app/tools/client"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
)

func RegisterRouter(e *echo.Group) {
	r := e.Group("/job")

	r.GET("", getList)
	r.GET("/:id", get)
	r.POST("", post)
	r.PUT("/:id", put)
	r.DELETE("/:id", _delete)
}

func getList(ctx echo.Context) error {
	u := ctx.Get("user").(*model.User)

	db.DB().Model(u).Where(&model.Job{Framework: ctx.QueryParam("framework")}).Related(&u.Jobs)
	//for i := range u.Jobs {
	//	job := &u.Jobs[i]
	//	job.Status = client.NewClient(client.NAMESPACE_DEEP_AI).GetJobStatus(fmt.Sprintf("user%d-job%d", u.ID, job.ID), job.Framework)
	//}
	db.DB().Model(u).Update(u)
	return ctx.JSON(http.StatusOK, u.Jobs)
}

func get(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))
	u := ctx.Get("user").(*model.User)
	job := new(model.Job)
	db.DB().Where(&model.Job{UserID: u.ID}).First(job, uint(id))
	//if job.Status == model.JOB_STATUS_TRAINING {
	//	job.Status = client.NewClient(client.NAMESPACE_DEEP_AI).GetJobStatus(fmt.Sprintf("user%d-job%d", u.ID, job.ID), job.Framework)
	//}
	return ctx.JSON(http.StatusOK, job)
}

func post(ctx echo.Context) error {
	u := ctx.Get("user").(*model.User)
	data := struct {
		Name             string
		Description      string
		Framework        string
		FrameworkVersion string
		EntryPoint       string
		Args             []string

		Num    int
		CPU    int
		Memory string
		GPU    int

		CodeID    uint
		DatasetID uint
		AIModelID uint
	}{}
	_ = ctx.Bind(&data)
	//if data.Framework != model.JOB_FRAMEWORK_TENSORFLOW && data.Framework != model.JOB_FRAMEWORK_PYTORCH {
	//	return echo.ErrBadRequest
	//}

	code, dataset, AIModel := new(model.Code), new(model.Dataset), new(model.AIModel)
	if db.DB().Where(&model.Code{UserID: u.ID}).First(code, data.CodeID).RecordNotFound() ||
		db.DB().Where(&model.Dataset{UserID: u.ID}).First(dataset, data.DatasetID).RecordNotFound() ||
		db.DB().Where(&model.AIModel{UserID: u.ID}).First(AIModel, data.AIModelID).RecordNotFound() {
		return echo.ErrBadRequest
	}

	job := model.Job{
		Name:             data.Name,
		Description:      data.Description,
		Framework:        data.Framework,
		FrameworkVersion: data.FrameworkVersion,
		EntryPoint:       data.EntryPoint,
		Args:             strings.Join(data.Args, " "),
		Num:              data.Num,
		CPU:              data.CPU, // TODO: CPU、内存、GPU资源限制
		Memory:           data.Memory,
		GPU:              data.GPU,
		UserID:           u.ID,
		CodeID:           data.CodeID,
		DatasetID:        data.DatasetID,
		AIModelID:        data.AIModelID,
	}
	db.DB().Create(&job)

	stopListener := make(chan struct{})
	go createJob(job, *code, *dataset, *AIModel, data.Args, stopListener)

	return ctx.JSON(http.StatusCreated, echo.Map{"id": job.ID})
}

func put(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))
	u := ctx.Get("user").(*model.User)

	if _, ok := ctx.QueryParams()["stop"]; ok {
		jobName := fmt.Sprintf("user%d-job%d", u.ID, id)
		job := new(model.Job)
		db.DB().Where(&model.Job{UserID: u.ID}).First(job, uint(id))

		cli := client.NewClient(client.NAMESPACE_DEEP_AI)
		if err := cli.DeleteTrainingJob(jobName, job.Framework); err != nil {
			return err
		}
		db.DB().Model(&model.Job{ID: uint(id), UserID: u.ID}).Update(model.Job{Status: model.JOB_STATUS_STOPPED})
		return ctx.NoContent(http.StatusOK)
	}

	data := struct {
		Name        string
		Description string
	}{}
	_ = ctx.Bind(&data)
	db.DB().Model(&model.Job{ID: uint(id), UserID: u.ID}).Update(model.Job{Name: data.Name, Description: data.Description})
	return ctx.NoContent(http.StatusOK)
}

func _delete(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))
	u := ctx.Get("user").(*model.User)
	jobName := fmt.Sprintf("user%d-job%d", u.ID, id)
	job := new(model.Job)
	db.DB().Where(&model.Job{UserID: u.ID}).First(job, uint(id))

	cli := client.NewClient(client.NAMESPACE_DEEP_AI)
	if err := cli.DeleteTrainingJob(jobName, job.Framework); err != nil {
		return err
	}
	db.DB().Delete(&model.Job{ID: uint(id), UserID: u.ID})
	return ctx.NoContent(http.StatusNoContent)
}
