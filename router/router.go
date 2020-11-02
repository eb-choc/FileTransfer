package router

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/kataras/iris"
)

type DirFile struct {
	IsDir bool
	Name  string
}

func InitRouter(app *iris.Application) {

	app.Get("/", func(ctx iris.Context) {
		loadDirFiles(ctx, "")
		ctx.ViewData("uploaded", ctx.URLParam("uploaded"))
		ctx.ViewData("time", ctx.URLParam("time"))
		ctx.View("index.html")
	})
	app.Get("/index", func(ctx iris.Context) {
		loadDirFiles(ctx, "")
		ctx.ViewData("uploaded", ctx.URLParam("uploaded"))
		ctx.ViewData("time", ctx.URLParam("time"))
		ctx.View("index.html")
	})
	app.Post("/files/transfer", func(ctx iris.Context) {
		file, info, err := ctx.FormFile("file")
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.HTML("Error while uploading: <b>" + err.Error() + "</b>")
			return
		}
		defer file.Close()
		out, err := os.OpenFile("./uploads/"+info.Filename, os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.HTML("Error while uploading: <b>" + err.Error() + "</b>")
			return
		}
		defer out.Close()
		io.Copy(out, file)
		ctx.Redirect("/index?uploaded=1&time=" + time.Now().Format("2006-01-02 15:04:05"));
	})
	app.Get("/dir", func(ctx iris.Context) {
		path := ctx.URLParam("path")
		file := ctx.URLParam("file")
		loadDirFiles(ctx, path)
		if file != "" {
			err := ctx.SendFile("." + path + "/" + file, file)
			if err != nil {
				log.Fatal(err.Error())
			}
		}
		ctx.View("index.html")
	})
}
func getCurrentDir() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	rst := filepath.Dir(path)
	return rst
}
func getImgBase64(file string) string{
	ff, _ := os.Open(file)
	defer ff.Close()
	sourcebuffer := make([]byte, 5000000)
	n, _ := ff.Read(sourcebuffer)
	//base64压缩
	sourcestring := base64.StdEncoding.EncodeToString(sourcebuffer[:n])
	return "data:image/png;base64," + sourcestring
}

func loadDirFiles(ctx iris.Context, relativePath string) {
	curDir := getCurrentDir()
	baseDir := curDir
	if relativePath != "" {
		curDir += relativePath
	}
	fileInfoList, err := ioutil.ReadDir(curDir)
	if err != nil {
		log.Fatal(err)
	}
	var dirs []DirFile
	var files []DirFile
	fmt.Println(len(fileInfoList))
	for _, v := range fileInfoList {
		if v.IsDir() {
			dirs = append(dirs, DirFile{
				Name:  v.Name(),
				IsDir: v.IsDir(),
			})
		} else {
			files = append(files, DirFile{
				Name:  v.Name(),
				IsDir: v.IsDir(),
			})
		}
	}
	for _, v := range files {
		dirs = append(dirs, v)
	}
	ctx.ViewData("dirs", dirs)

	ctx.ViewData("img_dir", getImgBase64(baseDir + "\\imgs\\dir.png"));
	ctx.ViewData("img_file", getImgBase64(baseDir + "\\imgs\\file.png"));
	ctx.ViewData("parentDir", relativePath)
}
