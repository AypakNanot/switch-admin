package main

import (
	"fmt"
	"log"
	"time"

	_ "github.com/GoAdminGroup/go-admin/adapter/gin"
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/sqlite"
	_ "github.com/GoAdminGroup/themes/sword"

	"github.com/GoAdminGroup/go-admin/engine"
	"switch-admin/internal/datamodel"
	"switch-admin/internal/handler"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/language"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/chartjs"
	"github.com/GoAdminGroup/themes/adminlte"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	// gin.DefaultWriter = ioutil.Discard  // 注释掉以启用日志

	// 主引擎用于 GoAdmin
	r := gin.Default()
	e := engine.Default()

	sysHandler := handler.NewSystemHandler()

	// 在 GoAdmin 之前注册 API 路由（直接注册到主路由器）
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, map[string]interface{}{"status": "ok"})
	})
	r.GET("/api/mode", sysHandler.GinAPIGetMode)
	r.POST("/api/mode", sysHandler.GinAPISwitchMode)
	r.GET("/api/system/config", func(c *gin.Context) {
		currentMode, err := sysHandler.GetRunMode()
		if err != nil {
			currentMode = "mock"
		}
		modeDesc := "离线测试模式"
		if currentMode == "switch" {
			modeDesc = "交换机模式"
		}
		c.JSON(200, map[string]interface{}{
			"code": 200,
			"data": map[string]interface{}{
				"mode":            currentMode,
				"mode_description": modeDesc,
				"database":        "SQLite3 (data/admin.db)",
				"goadmin_version": "v1.2.26",
			},
		})
	})
	// 首页重定向到 Dashboard
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/admin")
	})
	r.GET("/admin", func(c *gin.Context) {
		c.Redirect(302, "/admin/dashboard")
	})

	cfg := config.Config{
		Env: config.EnvLocal,
		Databases: config.DatabaseList{
			"default": {
				Driver:          config.DriverSqlite,
				File:            "data/admin.db",
				MaxIdleConns:    50,
				MaxOpenConns:    150,
				ConnMaxLifetime: time.Hour,
			},
		},
		UrlPrefix:           "admin",
		Store: config.Store{
			Path:   "./uploads",
			Prefix: "uploads",
		},
		Language:           language.CN,
		IndexUrl:           "/",
		Debug:              true,
		AccessAssetsLogOff: true,
		Animation: config.PageAnimation{
			Type: "fadeInUp",
		},
		ColorScheme:       adminlte.ColorschemeSkinBlack,
		BootstrapFilePath: "./internal/datamodel/bootstrap.go",
	}

	template.AddComp(chartjs.NewChart())

	if err := e.AddConfig(&cfg).
		AddGenerators(datamodel.Generators).
		AddGenerator("user", datamodel.GetUserTable).
		AddDisplayFilterXssJsFilter().
		Use(r); err != nil {
		panic(err)
	}

	r.Static("/uploads", "./uploads")

	// 注册自定义页面 - 使用 GoAdmin 的 HTML 方法
	e.HTML("GET", "/admin/dashboard", datamodel.GetDashboardContent, false)
	e.HTML("GET", "/admin/system/config", datamodel.GetSystemConfigPage, false)

	fmt.Println("=== GoAdmin 启动完成 ===")
	fmt.Println("Admin UI: http://localhost:9033/admin")
	fmt.Println("API: http://localhost:9033/api/mode")

	// 启动服务器
	if err := r.Run(":9033"); err != nil {
		log.Fatal("Server failed:", err)
	}

	log.Print("closing database connection")
	e.SqliteConnection().Close()
}
