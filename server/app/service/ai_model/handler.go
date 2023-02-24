package ai_model

import (
	"deep-ai-server/app/db"
	"deep-ai-server/app/model"
	"deep-ai-server/app/tools/client"
	"deep-ai-server/app/tools/file"
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func RegisterRouter(e *echo.Group) {
	r := e.Group("/model")

	r.GET("", getList)
	r.GET("/:id", get)
	r.GET("/:id/dir", getDir)
	r.GET("/:id/download", download)
	r.POST("", post)
	r.PUT("/:id", put)
	r.DELETE("/:id", _delete)
}

func getList(ctx echo.Context) error {
	u := ctx.Get("user").(*model.User)
	db.DB().Model(u).Related(&u.AIModels)

	if ctx.QueryParam("framework") == model.AIMODEL_FRAMEWORK_TENSORFLOW {
		db.DB().Model(u).Where(&model.AIModel{
			Framework: model.AIMODEL_FRAMEWORK_TENSORFLOW,
			Status: model.AIMODEL_STATUS_IDLE,
		}).Related(&u.AIModels)
	} else if ctx.QueryParam("framework") == model.AIMODEL_FRAMEWORK_PYTORCH {
		db.DB().Model(u).Where(&model.AIModel{
			Framework: model.AIMODEL_FRAMEWORK_PYTORCH,
			Status: model.AIMODEL_STATUS_IDLE,
		}).Related(&u.AIModels)
	} else {
		db.DB().Model(u).Related(&u.AIModels)
	}

	return ctx.JSON(http.StatusOK, u.AIModels)
}

func get(ctx echo.Context) error {
	u := ctx.Get("user").(*model.User)
	AIModelID, _ := strconv.Atoi(ctx.Param("id"))

	AIModel := new(model.AIModel)
	db.DB().Where(&model.AIModel{UserID: u.ID}).First(AIModel, AIModelID)

	return ctx.JSON(http.StatusOK, AIModel)
}

func getDir(ctx echo.Context) error {
	u := ctx.Get("user").(*model.User)
	AIModelID, _ := strconv.Atoi(ctx.Param("id"))
	path := ctx.QueryParam("path")
	AIModel := new(model.AIModel)
	db.DB().First(AIModel, model.AIModel{ID: uint(AIModelID), UserID: u.ID})
	path = filepath.Join(AIModel.PersistentVolumePath, path)
	//path, _ = filepath.Abs(path)
	// TODO:
	path = filepath.ToSlash(path)
	if !strings.HasPrefix(path, AIModel.PersistentVolumePath) {
		return echo.ErrBadRequest
	}

	data := make([]echo.Map, 0)
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		data = append(data, echo.Map{
			"name":  f.Name(),
			"isDir": f.IsDir(),
		})
	}

	return ctx.JSON(http.StatusOK, data)
}

func download(ctx echo.Context) error {
	u := ctx.Get("user").(*model.User)
	AIModelID, _ := strconv.Atoi(ctx.Param("id"))
	AIModel := new(model.AIModel)
	db.DB().Where(&model.AIModel{UserID: u.ID}).First(AIModel, AIModelID)

	zipPath := filepath.Join(AIModel.PersistentVolumePath, fmt.Sprintf("%s.zip", AIModel.Name))
	fmt.Println(zipPath)
	err := file.Zip(AIModel.PersistentVolumePath, zipPath)
	fmt.Println(err)
	defer func() { _ = os.Remove(zipPath) }()
	return ctx.File(zipPath)
}

func post(ctx echo.Context) error {
	u := ctx.Get("user").(*model.User)

	// 参数
	data := struct {
		Name             string
		Description      string
		Framework        string
		FrameworkVersion string
		Network          string
	}{}
	_ = ctx.Bind(&data)

	AIModel := model.AIModel{
		Name:             data.Name,
		Description:      data.Description,
		Framework:        data.Framework,
		FrameworkVersion: data.FrameworkVersion,
		Network:          data.Network,
		UserID:           u.ID,
	}
	db.DB().Create(&AIModel)

	c := make(chan string)
	stopPVCEventListener := make(chan struct{})
	go createPVC(&AIModel, c, stopPVCEventListener)
	go createNNModelFile(&AIModel, c)

	return ctx.JSON(http.StatusCreated, echo.Map{"id": AIModel.ID})
}

func put(ctx echo.Context) error {
	u := ctx.Get("user").(*model.User)
	AIModelID, _ := strconv.Atoi(ctx.Param("id"))

	data := struct {
		Name        string
		Description string
	}{}
	_ = ctx.Bind(&data)
	db.DB().Model(&model.AIModel{ID: uint(AIModelID), UserID: u.ID}).Update(model.AIModel{Name: data.Name, Description: data.Description})
	return ctx.NoContent(http.StatusOK)
}

func _delete(ctx echo.Context) error {
	u := ctx.Get("user").(*model.User)
	AIModelID, _ := strconv.Atoi(ctx.Param("id"))
	pvcName := fmt.Sprintf("user%d-model%d", u.ID, AIModelID)

	cli := client.NewClient(client.NAMESPACE_DEEP_AI)
	if err := cli.DeletePVC(pvcName); err != nil {
		return err
	}
	db.DB().Delete(&model.AIModel{ID: uint(AIModelID), UserID: u.ID})
	return ctx.NoContent(http.StatusNoContent)
}
