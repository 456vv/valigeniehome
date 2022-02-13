package valigeniehome

import (
	"sync"
	"net/http"
	"encoding/json"
	"fmt"
	"log"
	"io"
	"bytes"
	"github.com/tidwall/gjson"
)

// aligenie
type Aligenie struct {
	AppIdAttr 		string				// 属性id
	ErrorLog		*log.Logger			// 错误日志
	Debug			bool
	
	apps	map[string]*Application		// app集
	m		sync.Mutex					// 锁
}
func (T *Aligenie) logf(calldepth int, format string, a ...interface{}) string {
	txt := fmt.Sprintf(format+"\n", a...)
	if T.ErrorLog != nil {
		T.ErrorLog.Output(calldepth, txt)
	}
	log.Output(calldepth, txt)
	return txt
}
func (T *Aligenie) Logf(format string, a ...interface{}) string {
	return T.logf(2, format, a...)
}

//设置APP
//	id string			id名称
//	app *Application	app配置
func (T *Aligenie) SetApp(id string, app *Application){
	T.m.Lock()
	defer T.m.Unlock()
	if T.apps == nil {
		T.apps = make(map[string]*Application)
	}
	if app == nil {
		delete(T.apps, id)
		return
	}
	T.apps[id]	= app
}

//服务处理
//	w http.ResponseWriter	http响应
//	r *http.Request			http请求
func (T *Aligenie) ServeHTTP(w http.ResponseWriter, r *http.Request){
	T.m.Lock()
	var (
		query		= r.URL.Query()
		argId		= query.Get(T.AppIdAttr)
		req			Requester
		res 		Responser
	)
	
	if argId == "" {
		argId = r.Header.Get(T.AppIdAttr)
	}
	app, ok := T.apps[argId]
	T.m.Unlock()
	
	if !ok {
		logText := T.Logf("valigeniehome: 参数验证不通过,error(%s=%v)", T.AppIdAttr, argId)
		http.Error(w, logText,  http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	var buf io.Reader = r.Body
	
	b, err := io.ReadAll(buf)
	if err != nil {
		logText := T.Logf("valigeniehome: 读取Body数据失败,error(%v)", err)
		http.Error(w, logText,  http.StatusBadRequest)
		return
	}
	
	if T.Debug {
		T.Logf("valigeniehome: 请求>%s", b)
	}
	
	if !json.Valid(b) {
		logText := T.Logf("valigeniehome: 主体不是有效的JSON格式")
		http.Error(w, logText,  http.StatusBadRequest)
		return
	}
	
	payLoadVersion := gjson.GetBytes(b, "header.payLoadVersion")
	if payLoadVersion.Type != gjson.Number {
		logText := T.Logf("valigeniehome: JSON格式无法识别")
		http.Error(w, logText,  http.StatusBadRequest)
		return
	}
	
	buf = bytes.NewBuffer(b)
	switch ver := payLoadVersion.Uint(); ver {
		case 1:{
			reqs := &reqSynthesis{v1:new(Request1)}
			ress := &resSynthesis{v1:new(Response1)}
			ress.v1.aligenie = T
			ress.v1.req = reqs.v1
			ress.version = 1
			req = reqs
			res = ress
			err = json.NewDecoder(buf).Decode(reqs.v1)
		}
		case 2:{
			reqs := &reqSynthesis{v2:new(Request2)}
			ress := &resSynthesis{v2:new(Response2)}
			ress.v2.aligenie = T
			ress.v2.req = reqs.v2
			ress.version = 2
			req = reqs
			res = ress
			err = json.NewDecoder(buf).Decode(reqs.v2)
		}
		default:{
			logTxt := T.Logf("valigeniehome: 暂不支持(%d)版本，请联系厂商更新！\n", ver)
			http.Error(w, logTxt, http.StatusBadRequest)
			return
		}
	}
	
	if err != nil {
		logText := T.Logf("valigeniehome: 数据验证不通过,error(%s)\n请求BODY：%s\n\n", err.Error(), b)
		http.Error(w, logText,  http.StatusBadRequest)
		return
	}
	
	home := Home{Request: req, Response: res, app: app}
	home.serveHTTP(w, r)
}













