package valigeniehome


import (
	"net/http"
)

// 程序
type Application struct{
	HandleFunc				http.HandlerFunc							// 处理函数
}


