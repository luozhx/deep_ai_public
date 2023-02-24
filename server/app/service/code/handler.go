package code

import (
	"deep-ai-server/app/db"
	"deep-ai-server/app/model"
	"deep-ai-server/app/tools/client"
	"deep-ai-server/app/tools/file"
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func RegisterRouter(e *echo.Group) {
	r := e.Group("/code")

	r.GET("", getList)
	r.GET("/:id", get)
	r.GET("/:id/dir", getDir)
	r.GET("/:id/download", download)
	r.POST("", post)
	r.PUT("/:id", put)
	r.DELETE("/:id", _delete)

	r.Any("/:id/notebook/:service_name/*", proxy)
}

func getList(ctx echo.Context) error {
	u := ctx.Get("user").(*model.User)
	db.DB().Model(u).Where(&model.Code{
		Framework: ctx.QueryParam("framework"),
	}).Related(&u.Codes)
	return ctx.JSON(http.StatusOK, u.Codes)
}

func get(ctx echo.Context) error {
	u := ctx.Get("user").(*model.User)
	codeID, _ := strconv.Atoi(ctx.Param("id"))

	code := new(model.Code)
	db.DB().Where(&model.Code{UserID: u.ID}).First(code, codeID)

	return ctx.JSON(http.StatusOK, code)
}

func getDir(ctx echo.Context) error {
	u := ctx.Get("user").(*model.User)
	codeID, _ := strconv.Atoi(ctx.Param("id"))
	path := ctx.QueryParam("path")
	code := new(model.Code)
	db.DB().First(code, model.Code{ID: uint(codeID), UserID: u.ID})
	path = filepath.Join(code.PersistentVolumePath, path)
	//path, _ = filepath.Abs(path)
	// TODO:
	path = filepath.ToSlash(path)
	if !strings.HasPrefix(path, code.PersistentVolumePath) {
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
	codeID, _ := strconv.Atoi(ctx.Param("id"))
	code := new(model.Code)
	db.DB().Where(&model.Code{UserID: u.ID}).First(code, codeID)

	if code.PVCStatus != model.PVC_STATUS_BOUND {
		return echo.ErrNotFound // TODO: PVC未绑定时，不允许下载
	}

	zipPath := filepath.Join(code.PersistentVolumePath, fmt.Sprintf("%s.zip", code.Name))
	_ = file.Zip(code.PersistentVolumePath, zipPath)
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
	}{}
	_ = ctx.Bind(&data)
	code := model.Code{
		Name:             data.Name,
		Description:      data.Description,
		Framework:        data.Framework,
		FrameworkVersion: data.FrameworkVersion,
		UserID:           u.ID,
	}
	db.DB().Create(&code)

	c := make(chan string)
	stopPVCEventListener := make(chan struct{})
	stopDeploymentEventListener := make(chan struct{})

	go createPVC(code, c, stopPVCEventListener)
	go createNotebook(code, c, stopDeploymentEventListener)

	return ctx.JSON(http.StatusCreated, echo.Map{"id": code.ID})
}

func put(ctx echo.Context) error {
	u := ctx.Get("user").(*model.User)
	codeID, _ := strconv.Atoi(ctx.Param("id"))

	data := struct {
		Name        string
		Description string
	}{}
	_ = ctx.Bind(&data)

	db.DB().Model(&model.Code{ID: uint(codeID), UserID: u.ID}).Update(model.Code{Name: data.Name, Description: data.Description})

	return ctx.NoContent(http.StatusOK)
}

func _delete(ctx echo.Context) error {
	u := ctx.Get("user").(*model.User)
	codeID, _ := strconv.Atoi(ctx.Param("id"))
	pvcName := fmt.Sprintf("user%d-code%d", u.ID, codeID)

	cli := client.NewClient(client.NAMESPACE_DEEP_AI)
	if err := cli.DeletePVC(pvcName); err != nil {
		return err
	}
	if err := cli.DeleteNotebook(pvcName); err != nil { // TODO: NOTEBOOK NOT FOUND 忽略错误
		return err
	}
	db.DB().Delete(&model.Code{ID: uint(codeID), UserID: u.ID}) // 由用户ID与代码ID删除对应的代码
	return ctx.NoContent(http.StatusNoContent)
}

// reference: https://stackoverflow.com/questions/42468552/how-to-handle-proxying-to-multiple-services-with-golang-and-labstack-echo
func proxy(c echo.Context) error {
	//todo: verify user
	name := c.Param("service_name")
	handler := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s.deep-ai.svc.cluster.local", name),
	})
	handler.ServeHTTP(c.Response(), c.Request()) //	reference: echo.WrapHandler()
	return nil
}
