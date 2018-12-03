package valigeniehome

import (
	"net/http"
	"context"
)

// 家
type Home struct{
    Response		*Response											// 响应
    Request     	*Request											// 请求
    App        	 	*Application										// app
}


//服务处理
//	w http.ResponseWriter	http响应对象
// 	r *http.Request			http请求对象
func (T *Home) ServeHTTP(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	if T.App.HandleFunc != nil {
		r = r.WithContext(context.WithValue(r.Context(), r.URL.Path, T))
		T.App.HandleFunc(w, r)
		return
	}
	http.Error(w, "你没有设置 valigeniohome.Application.HandleFunc",  http.StatusInternalServerError)
}