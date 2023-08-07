package main

import (
	"go-app/router"
)

func main() {
	// 載入路由
	r := router.InitRouter()
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
