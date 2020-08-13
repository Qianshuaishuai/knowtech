package main

import (
	"github.com/astaxie/beego"
	"fmt"
	"runtime"
	"knowtech/models"

	"github.com/astaxie/beego/context"

	_ "knowtech/routers"
	"knowtech/controllers"
	"github.com/HYY-yu/LogLib"
	"errors"
)

func main() {
	//连接数据库
	models.InitGorm()
	db := models.GetDb()
	defer db.Close()

	if !models.CheckDB(db){
		fmt.Println("数据库错误，无法启动")
		return
	}

	//如果服务器Panic ，返回500错，而不是错误信息。
	beego.BConfig.RecoverFunc = recoverFuncForMyServer
	beego.ErrorController(&controllers.ErrorController{})
	beego.Run()
}

func recoverFuncForMyServer(ctx *context.Context) {
	if err := recover(); err != nil {
		if err == beego.ErrAbort {
			return
		}
		if !models.MyConfig.RecoverPanic {
			panic(err)
		}

		//通知客户端
		beego.Exception(500, ctx)

		var stack []string
		stack = make([]string, 0, 4)
		for i := 1; ; i++ {
			if _, file, line, ok := runtime.Caller(i); ok {
				stack = append(stack, fmt.Sprintf("%s:%d", file, line))
			} else {
				break
			}
		}

		//记录此错误
		realErr, ok := err.(error)
		if !ok {
			realErr = errors.New(fmt.Sprint(err))
		}
		loglib.GetLogger().LogPanicForServer(realErr, stack, ctx)
	}
}
