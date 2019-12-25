package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var mux map[string]func(http.ResponseWriter, *http.Request)

type Myhandler struct{}
type home struct {
	Title string
}

const (
	Template_Dir = "./view/"
	Upload_Dir   = "./upload/"
)

func main() {
	server := http.Server{
		Addr:        ":8080",
		Handler:     &Myhandler{},
		ReadTimeout: 10 * time.Second,
	}
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/"] = index
	mux["/clear"] = clear
	mux["/upload"] = upload
	mux["/file"] = StaticServer
	server.ListenAndServe()
}

func (*Myhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)
		return
	}
	http.StripPrefix("/", http.FileServer(http.Dir("./upload/"))).ServeHTTP(w, r)
}

func upload(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		t, _ := template.ParseFiles(Template_Dir + "index.html")
		t.Execute(w, "上传文件")
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Fprintf(w, "%v", "上传错误")
			return
		}
		fileext := filepath.Ext(handler.Filename)
		if check(fileext) == false {
			fmt.Fprintf(w, "%v", "不允许的上传类型")
			return
		}
		fmt.Println(handler.Filename)
		// strconv.FormatInt(time.Now().Unix(), 10) + "." +
		urifilename := handler.Filename
		f, err := os.OpenFile(Upload_Dir+urifilename, os.O_CREATE|os.O_WRONLY, 0660)
		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, "%v", "文件创建错误")
			return
		}
		_, err = io.Copy(f, file)
		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, "%v", "上传失败")
			return
		}
		//filedir, _ := filepath.Abs(Upload_Dir + filename)
		fmt.Fprintf(w, "%v", "上传完成,服务器地址:"+urifilename)
	}
}

func clear(w http.ResponseWriter, r *http.Request) {
	os.RemoveAll(Upload_Dir)
	os.Mkdir(Upload_Dir, os.ModePerm)
	fmt.Fprintf(w, "%v", "清理成功")
}

func index(w http.ResponseWriter, r *http.Request) {
	title := home{Title: "首页"}
	t, _ := template.ParseFiles(Template_Dir + "index.html")
	t.Execute(w, title)
}

func StaticServer(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/file", http.FileServer(http.Dir("./upload/"))).ServeHTTP(w, r)
}

func check(name string) bool {
	ext := []string{".exe", ".js", ".png"}

	for _, v := range ext {
		if v == name {
			return false
		}
	}
	return true
}
