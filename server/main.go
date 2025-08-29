package main

import (
	"cooller/server/api/v1/business"
	"cooller/server/core"
	"cooller/server/global"
	"cooller/server/initialize"
	"fmt"
	"go.uber.org/zap"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

// @title                       Swagger Example API
// @version                     0.0.1
// @description                 This is a sample Server pets
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        x-token
// @BasePath                    /
func main() {
	global.GVA_VP = core.Viper() // 初始化Viper
	initialize.OtherInit()       // 本地缓存jwt过期时间等
	global.GVA_LOG = core.Zap()  // 初始化zap日志库
	zap.ReplaceGlobals(global.GVA_LOG)
	global.GVA_DB = initialize.Gorm() // gorm连接数据库
	initialize.Timer()
	_, err := global.GVA_Timer.AddTaskByFunc("CheckQrCode", "@every 168h", func() {
		business.CheckQrCodeExpired()
	})
	if err != nil {
		fmt.Println("add timer error:", err)
	}
	initialize.DBList()

	//tm := timer.NewTimer()
	//tm.ScheduleFunc(30*time.Minute, func() {
	//	business.CheckQrCodeExpired()
	//})
	//tm.Run()
	if global.GVA_DB != nil {
		initialize.RegisterTables() // 初始化表
		// 程序结束前关闭数据库链接
		db, _ := global.GVA_DB.DB()
		defer db.Close()
	}
	core.RunWindowsServer()
}
