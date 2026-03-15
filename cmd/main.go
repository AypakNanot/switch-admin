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
	routeHandler := handler.NewRouteHandler()
	diagnosticHandler := handler.NewDiagnosticHandler()

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

	// 网络模块 API - 路由管理
	r.GET("/api/v1/routes/table", routeHandler.GetRouteTable)
	r.GET("/api/v1/routes/static", routeHandler.GetStaticRoutes)
	r.GET("/api/v1/routes/static/:id", routeHandler.GetStaticRoute)
	r.POST("/api/v1/routes/static", routeHandler.CreateStaticRoute)
	r.PUT("/api/v1/routes/static/:id", routeHandler.UpdateStaticRoute)
	r.DELETE("/api/v1/routes/static/:id", routeHandler.DeleteStaticRoute)

	// 网络模块 API - 诊断工具
	r.GET("/api/v1/diagnostic/cable/ports", diagnosticHandler.GetDetectablePorts)
	r.POST("/api/v1/diagnostic/cable", diagnosticHandler.ExecuteCableTest)
	r.GET("/api/v1/diagnostic/ping/:task_id", diagnosticHandler.GetPingTaskResult)
	r.POST("/api/v1/diagnostic/ping", diagnosticHandler.CreatePingTask)
	r.DELETE("/api/v1/diagnostic/ping/:task_id", diagnosticHandler.DeletePingTask)
	r.GET("/api/v1/diagnostic/traceroute/:task_id", diagnosticHandler.GetTracerouteTaskResult)
	r.POST("/api/v1/diagnostic/traceroute", diagnosticHandler.CreateTracerouteTask)
	r.DELETE("/api/v1/diagnostic/traceroute/:task_id", diagnosticHandler.DeleteTracerouteTask)
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

	// 先添加配置（这会初始化数据库连接）
	e.AddConfig(&cfg)

	// 在 Use 之前先初始化数据库表（创建 GoAdmin 所需的表）
	datamodel.InitDatabaseTables(e.SqliteConnection())

	if err := e.AddGenerators(datamodel.Generators).
		AddGenerator("user", datamodel.GetUserTable).
		AddDisplayFilterXssJsFilter().
		Use(r); err != nil {
		panic(err)
	}

	// 初始化菜单
	datamodel.InitMenu(e.SqliteConnection())
	datamodel.InitDashboard(e.SqliteConnection())

	r.Static("/uploads", "./uploads")

	// 注册自定义页面 - 使用 GoAdmin 的 HTML 方法
	// Dashboard 和系统配置
	e.HTML("GET", "/admin/dashboard", datamodel.GetDashboardContent, false)
	e.HTML("GET", "/admin/system/config", datamodel.GetSystemConfigPage, false)

	// 网络模块 - IP 路由
	e.HTML("GET", "/admin/network/route-table", datamodel.GetRouteTableContent, false)
	e.HTML("GET", "/admin/network/static-route", datamodel.GetStaticRouteContent, false)

	// 网络模块 - 诊断工具
	e.HTML("GET", "/admin/network/ping", datamodel.GetPingContent, false)
	e.HTML("GET", "/admin/network/traceroute", datamodel.GetTracerouteContent, false)
	e.HTML("GET", "/admin/network/cable-test", datamodel.GetCableTestContent, false)

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
