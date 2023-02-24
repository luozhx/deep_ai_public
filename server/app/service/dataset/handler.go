package dataset

import (
	"deep-ai-server/app/db"
	"deep-ai-server/app/model"
	"deep-ai-server/app/tools/client"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"path/filepath"
	"strconv"
)

func RegisterRouter(e *echo.Group) {
	r := e.Group("/dataset")

	r.GET("", getList)
	r.GET("/:id", get)
	r.GET("/:id/download", download)
	r.POST("", post)
	r.PUT("/:id", put)
	r.DELETE("/:id", _delete)
}

func getList(ctx echo.Context) error {
	u := ctx.Get("user").(*model.User)
	db.DB().Model(u).Related(&u.Datasets)
	return ctx.JSON(http.StatusOK, u.Datasets)
}

func get(ctx echo.Context) error {
	u := ctx.Get("user").(*model.User)
	datasetID, _ := strconv.Atoi(ctx.Param("id"))

	dataset := new(model.Dataset)
	db.DB().Where(&model.Dataset{UserID: u.ID}).First(dataset, datasetID)

	return ctx.JSON(http.StatusOK, dataset)
}

func download(ctx echo.Context) error {
	u := ctx.Get("user").(*model.User)
	id, _ := strconv.Atoi(ctx.Param("id"))
	dataset := new(model.Dataset)
	db.DB().Where(&model.Dataset{UserID: u.ID}).First(dataset, id)

	zipPath := filepath.Join(dataset.PersistentVolumePath, fmt.Sprintf("%s.zip", dataset.Name))
	//_ = file.Zip(dataset.PersistentVolumePath, zipPath)
	//defer func() { _ = os.Remove(zipPath) }()
	return ctx.File(zipPath)
}

func post(ctx echo.Context) error {
	u := ctx.Get("user").(*model.User)

	// 参数
	data := struct {
		Name        string
		Description string
	}{}
	_ = ctx.Bind(&data)
	datasetZipFile, err := ctx.FormFile("file")

	if err != nil {
		return echo.ErrBadRequest
	}

	dataset := model.Dataset{
		Name:        data.Name,
		Description: data.Description,
		UserID:      u.ID,
		Size:        datasetZipFile.Size,
	}
	db.DB().Create(&dataset)

	c := make(chan string)
	stopPVCEventListener := make(chan struct{})
	go createPVC(&dataset, c, stopPVCEventListener)
	go saveDataset(&dataset, datasetZipFile, c)

	return ctx.JSON(http.StatusCreated, echo.Map{"id": dataset.ID})
}

func put(ctx echo.Context) error {
	u := ctx.Get("user").(*model.User)
	datasetID, _ := strconv.Atoi(ctx.Param("id"))

	data := struct {
		Name        string
		Description string
	}{}
	_ = ctx.Bind(&data)
	db.DB().Model(&model.Dataset{ID: uint(datasetID), UserID: u.ID}).Update(&model.Dataset{Name: data.Name, Description: data.Description})

	return ctx.NoContent(http.StatusOK)
}

func _delete(ctx echo.Context) error {
	u := ctx.Get("user").(*model.User)
	datasetID, _ := strconv.Atoi(ctx.Param("id"))

	pvcName := fmt.Sprintf("user%d-dataset%d", u.ID, datasetID)
	cli := client.NewClient(client.NAMESPACE_DEEP_AI)
	if err := cli.DeletePVC(pvcName); err != nil {
		return err
	}
	db.DB().Delete(&model.Dataset{ID: uint(datasetID), UserID: u.ID})
	return ctx.NoContent(http.StatusNoContent)
}
