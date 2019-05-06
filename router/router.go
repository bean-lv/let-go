package router

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"letgo/context"
	"letgo/controller"
	"letgo/plugins/cors"
	"net/http"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Router interface {
	ServeHTTP(rw http.ResponseWriter, req *http.Request)
	AddAutoRouter(c controller.Controller)
}

type myRouter struct {
	routerMap    map[string]route
	staticFolder string
	project      string
	homepage     string

	pool sync.Pool

	cors cors.CORS // 跨域访问
}

type route struct {
	pattern         string // 路由格式：/api/account/login
	controllerType  reflect.Type
	methodName      string
	methodInputType reflect.Type // 方法参数类型，转发路由时转换json数据
}

func NewRouter() Router {
	r := &myRouter{
		routerMap: make(map[string]route),
	}
	r.pool.New = func() interface{} {
		return context.New()
	}

	r.cors = cors.New()
	r.staticFolder = Static_Folder
	r.project = Project_Name
	r.homepage = Homepage

	return r
}

func (r *myRouter) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ctx := r.pool.Get().(context.Context)
	defer r.pool.Put(ctx)
	ctx.Reset(rw, req)

	// 判断是否支持跨域访问
	if r.cors != nil {
		r.cors.PrepareCors(ctx.Response(), ctx.Request())
	}

	if req.URL.Path == "/" {
		// 默认首页
		r.serveHtml(rw, req, path.Join(r.staticFolder, r.project, r.homepage))
	} else if strings.HasPrefix(req.URL.Path, Prefix_API) {
		// 处理API访问逻辑
		r.serveAPI(rw, req, ctx)
	} else if strings.HasPrefix(req.URL.Path, Prefix_Static) {
		// 静态资源
		r.serveFile(rw, req)
	} else if strings.HasPrefix(req.URL.Path, Prefix_Upload) {
		// 上传文件
		r.serveUpload(rw, req)
	}
}

func (r *myRouter) serveUpload(rw http.ResponseWriter, req *http.Request) {
	r.uploadOneFile(rw, req)
	// r.uploadMultiFiles(rw, req)
}

// func (r *RouterRegister) uploadMultiFiles(rw http.ResponseWriter, req *http.Request) {
// 	req.ParseMultipartForm(32 << 20)
// 	multipartForm := req.MultipartForm
// 	files := multipartForm.File["uploadfile"]

// 	var filenames []string
// 	for _, file := range files {

// 		src, err := file.Open()
// 		if err != nil {
// 			Error("Upload-> open file error:", err.Error())
// 			return
// 		}
// 		defer src.Close()

// 		filename, err := r.getUploadFileName(req.URL.Path, file.Filename)
// 		if err != nil {
// 			Error("Upload-> get file name error:", err.Error())
// 			return
// 		}
// 		dst, err := os.Create(filename)
// 		if err != nil {
// 			Error("Upload-> create file", filename, "error:", err.Error())
// 			return
// 		}
// 		defer dst.Close()

// 		if _, err = io.Copy(dst, src); err != nil {
// 			Error("Upload-> Copy file", filename, "error:", err.Error())
// 			return
// 		}
// 		filenames = append(filenames, filename)

// 	}

// 	json.NewEncoder(rw).Encode(filenames)
// }

func (r *myRouter) uploadOneFile(rw http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(32 << 20)

	file, handler, err := req.FormFile("uploadfile")
	if err != nil {
		return
	}
	defer file.Close()

	uploadFilename, err := r.getUploadFileName(req.URL.Path, handler.Filename)
	if err != nil {
		return
	}

	f, err := os.OpenFile(uploadFilename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		return
	}

	json.NewEncoder(rw).Encode(uploadFilename)
}

func (r *myRouter) getUploadFileName(url, filename string) (uploadFilename string, err error) {
	if len(filename) == 0 {
		err = errors.New("文件名不能为空")
	}
	url = url[1:]

	var folders []string
	if strings.Contains(url, "/") {
		folders = strings.Split(url, "/")
	} else if strings.Contains(url, "\\") {
		folders = strings.Split(url, "\\")
	} else {
		folders = []string{url}
	}
	p := path.Join(folders...)
	p = path.Join(r.staticFolder, p)

	var infos []string
	if strings.Contains(filename, ".") {
		infos = strings.Split(filename, ".")
	} else {
		infos = []string{filename}
	}
	if len(infos) == 1 {
		err = errors.New("不支持的文件格式")
		return
	}

	h := md5.New()
	io.WriteString(h, strconv.FormatInt(time.Now().Unix(), 10))
	io.WriteString(h, infos[0])
	f := fmt.Sprintf("%x.%s", h.Sum(nil), infos[len(infos)-1])

	uploadFilename = path.Join(p, f)
	return
}

func (r *myRouter) serveHtml(rw http.ResponseWriter, req *http.Request, url string) {
	t, err := template.ParseFiles(url)
	if err != nil {
		http.NotFound(rw, req)
		return
	}

	t.Execute(rw, nil)
}

func (r *myRouter) serveFile(rw http.ResponseWriter, req *http.Request) {
	var url = req.URL.Path
	url = url[1:]
	http.ServeFile(rw, req, url)
}

func (r *myRouter) serveAPI(rw http.ResponseWriter, req *http.Request, ctx context.Context) {
	route, err := r.findRouterInfo(req.URL.Path)
	if err != nil {
		http.NotFound(rw, req)
		return
	}

	var execController controller.Controller
	refV := reflect.New(route.controllerType)
	execController, ok := refV.Interface().(controller.Controller)
	if !ok {
		http.NotFound(rw, req)
		return
	}
	execController.Init(ctx)

	vc := reflect.ValueOf(execController)
	method := vc.MethodByName(route.methodName)
	if !method.IsValid() {
		http.NotFound(rw, req)
		return
	}

	methodInput, err := r.getMethodStructParams(route, req)
	if err != nil {
		json.NewEncoder(ctx.Response()).Encode(err.Error())
		return
	}

	inputs := []reflect.Value{}
	if methodInput != nil {
		expectType := method.Type().In(0)
		provideType := reflect.TypeOf(methodInput)
		if expectType != provideType {
			return
		}
		inputs = append(inputs, reflect.ValueOf(methodInput))
	}

	method.Call(inputs)
}

// 通过路由信息，转换请求数据为Struct类型的参数
func (r *myRouter) getMethodStructParams(route route, req *http.Request) (interface{}, error) {
	if route.methodInputType == nil {
		return nil, nil
	}

	body, err := ioutil.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		return nil, err
	}

	mit := reflect.New(route.methodInputType)
	err = json.Unmarshal(body, mit.Interface())
	if err != nil {
		return nil, err
	}
	return mit.Elem().Interface(), nil
}

func (r *myRouter) findRouterInfo(url string) (route route, err error) {
	// 格式：/api/account/login
	pattern := strings.ToLower(url)

	if routeInfo, ok := r.routerMap[pattern]; ok {
		route = routeInfo
		return
	}

	err = errors.New("Route not found")
	return
}

func (r *myRouter) AddAutoRouter(c controller.Controller) {
	reflectVal := reflect.ValueOf(c)
	rt := reflectVal.Type()
	ct := reflect.Indirect(reflectVal).Type()
	controllerName := strings.TrimSuffix(ct.Name(), Suffix_Controller)

	for i := 0; i < rt.NumMethod(); i++ {
		route := route{}
		route.controllerType = ct
		route.methodName = rt.Method(i).Name
		pattern := path.Join(Prefix_API, strings.ToLower(controllerName), strings.ToLower(rt.Method(i).Name))
		route.pattern = pattern

		// 提取方法的第一个参数信息，如果是Struct，则保存到路由信息，用户访问时json数据转换为Struct
		mt := rt.Method(i).Type

		if mt.NumIn() >= 2 {
			it := mt.In(1)
			if it.Kind() == reflect.Struct {
				route.methodInputType = it
			}
		}

		r.routerMap[pattern] = route
	}
}
