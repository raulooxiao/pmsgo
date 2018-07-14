package main

import (
	"fmt"
	"log"
	"myproject/pmsgo/models"
	"myproject/pmsgo/pkg/logging"
	"myproject/pmsgo/pkg/setting"
	"myproject/pmsgo/routers"
	"net/http"
	"time"
)

func main() {

	addr := fmt.Sprintf("%s:%d", setting.HTTPIP, setting.HTTPPort)

	router := routers.InitRouter()

	s := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    setting.ReadTimout * time.Second,
		WriteTimeout:   setting.WriteTimout * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// 定时更新密码表的记录
	go func() {
		for {
			timer := time.NewTicker(5 * time.Minute)
			select {
			case <-timer.C:
				models.UpdateIPPasswd()
			}
		}
	}()

	logging.Info("Now start PMS api server, listening: ", addr)

	// 启动api服务器
	err := s.ListenAndServe()
	if err != nil {
		log.Fatalf("Fail to start api svr: %v", err)
	}

}
