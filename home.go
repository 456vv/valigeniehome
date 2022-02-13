package valigeniehome

import (
	"net/http"
	"context"
)

//家
type Home struct{
    Request     	Requester											// 请求
    Response		Responser											// 响应
    app        	 	*Application										// app
}

//服务处理
//	w http.ResponseWriter	http响应对象
// 	r *http.Request			http请求对象
func (T *Home) serveHTTP(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	if T.app.HandleFunc != nil {
		r = r.WithContext(context.WithValue(r.Context(), HomeContextKey, T))
		T.app.HandleFunc(w, r)
		return
	}
	http.Error(w, "你没有设置 valigeniohome.Application.HandleFunc",  http.StatusInternalServerError)
}

