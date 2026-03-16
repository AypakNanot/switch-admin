package main

import (
	"fmt"
	"log"
	"time"

	_ "github.com/GoAdminGroup/go-admin/adapter/gin"
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/sqlite"
	_ "github.com/GoAdminGroup/themes/sword"
	_ "modernc.org/sqlite"

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
	maintenanceHandler := handler.NewMaintenanceHandler()
	configHandler := handler.NewConfigHandler()

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

	// 维护模块 API - 重启/保存
	r.POST("/api/v1/system/save-config", maintenanceHandler.SaveConfig)
	r.POST("/api/v1/system/reboot", maintenanceHandler.RebootSwitch)
	r.POST("/api/v1/system/factory-reset", maintenanceHandler.FactoryReset)

	// 维护模块 API - 系统配置
	r.GET("/api/v1/system/config", maintenanceHandler.GetSystemConfig)
	r.PUT("/api/v1/system/network", maintenanceHandler.UpdateNetworkConfig)
	r.PUT("/api/v1/system/temperature", maintenanceHandler.UpdateTemperatureConfig)
	r.PUT("/api/v1/system/info", maintenanceHandler.UpdateDeviceInfo)
	r.PUT("/api/v1/system/datetime", maintenanceHandler.UpdateDateTime)

	// 维护模块 API - 加载配置
	r.GET("/api/v1/config/files", maintenanceHandler.GetConfigFiles)
	r.POST("/api/v1/config/load", maintenanceHandler.LoadConfig)

	// 维护模块 API - 文件管理
	r.GET("/api/v1/files", maintenanceHandler.GetFiles)
	r.POST("/api/v1/files/upload", maintenanceHandler.UploadFile)
	r.POST("/api/v1/files/firmware", maintenanceHandler.UploadFirmware)
	r.GET("/api/v1/files/download", maintenanceHandler.DownloadFile)
	r.DELETE("/api/v1/files", maintenanceHandler.DeleteFiles)

	// 维护模块 API - 日志管理
	r.GET("/api/v1/logs", maintenanceHandler.GetLogs)
	r.DELETE("/api/v1/logs", maintenanceHandler.ClearLogs)

	// 维护模块 API - SNMP 配置
	r.GET("/api/v1/snmp/config", maintenanceHandler.GetSNMPConfig)
	r.PUT("/api/v1/snmp/config", maintenanceHandler.UpdateSNMPConfig)
	r.GET("/api/v1/snmp/communities", maintenanceHandler.GetSNMPCommunity)
	r.POST("/api/v1/snmp/communities", maintenanceHandler.AddSNMPCommunity)
	r.DELETE("/api/v1/snmp/communities/:name", maintenanceHandler.DeleteSNMPCommunity)

	// 维护模块 API - SNMP Trap 配置
	r.GET("/api/v1/snmp/trap/config", maintenanceHandler.GetSNMPTrapConfig)
	r.PUT("/api/v1/snmp/trap/config", maintenanceHandler.UpdateSNMPTrapConfig)
	r.GET("/api/v1/snmp/trap/hosts", maintenanceHandler.GetSNMPTrapHosts)
	r.POST("/api/v1/snmp/trap/hosts", maintenanceHandler.AddSNMPTrapHost)
	r.DELETE("/api/v1/snmp/trap/hosts/:id", maintenanceHandler.DeleteSNMPTrapHost)
	r.POST("/api/v1/snmp/trap/hosts/:id/test", maintenanceHandler.TestSNMPTrap)

	// 维护模块 API - 安全防护 - 蠕虫攻击防护
	r.GET("/api/v1/security/worm/rules", maintenanceHandler.GetWormRules)
	r.POST("/api/v1/security/worm/rules", maintenanceHandler.AddWormRule)
	r.PUT("/api/v1/security/worm/rules/:id", maintenanceHandler.UpdateWormRule)
	r.DELETE("/api/v1/security/worm/rules", maintenanceHandler.DeleteWormRules)
	r.POST("/api/v1/security/worm/clear-stats", maintenanceHandler.ClearWormStats)

	// 维护模块 API - 安全防护 - DDoS 攻击防护
	r.GET("/api/v1/security/ddos/config", maintenanceHandler.GetDDoSConfig)
	r.PUT("/api/v1/security/ddos/config", maintenanceHandler.UpdateDDoSConfig)

	// 维护模块 API - 安全防护 - ARP 攻击防护
	r.GET("/api/v1/security/arp/config", maintenanceHandler.GetARPConfig)
	r.PUT("/api/v1/security/arp/config", maintenanceHandler.UpdateARPConfig)

	// 维护模块 API - 用户管理
	r.GET("/api/v1/users", maintenanceHandler.GetUsers)
	r.POST("/api/v1/users", maintenanceHandler.CreateUser)
	r.PUT("/api/v1/users/:username", maintenanceHandler.UpdateUser)
	r.DELETE("/api/v1/users", maintenanceHandler.DeleteUsers)

	// 维护模块 API - 当前会话
	r.GET("/api/v1/sessions", maintenanceHandler.GetSessions)
	r.DELETE("/api/v1/sessions/:session_id", maintenanceHandler.DeleteSession)

	// 配置模块 API - 端口管理
	r.GET("/api/v1/ports", configHandler.GetPorts)
	r.PUT("/api/v1/ports/:port_id", configHandler.UpdatePort)

	// 配置模块 API - 链路聚合
	r.GET("/api/v1/link-aggregation", configHandler.GetLinkAggregation)
	r.POST("/api/v1/link-aggregation", configHandler.CreateLinkAggregation)
	r.PUT("/api/v1/link-aggregation/:id", configHandler.UpdateLinkAggregation)
	r.DELETE("/api/v1/link-aggregation/:id", configHandler.DeleteLinkAggregation)

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
	datamodel.InitConfigMenu(e.SqliteConnection())

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

	// 维护模块
	e.HTML("GET", "/admin/maintenance/reboot-save", datamodel.GetRebootSaveContent, false)
	e.HTML("GET", "/admin/maintenance/users", datamodel.GetUsersContent, false)
	e.HTML("GET", "/admin/maintenance/system-config", datamodel.GetMaintenanceSystemConfigContent, false)
	e.HTML("GET", "/admin/maintenance/load-config", datamodel.GetLoadConfigContent, false)
	e.HTML("GET", "/admin/maintenance/files", datamodel.GetFilesContent, false)
	e.HTML("GET", "/admin/maintenance/logs", datamodel.GetLogsContent, false)
	e.HTML("GET", "/admin/maintenance/snmp", datamodel.GetSNMPContent, false)
	e.HTML("GET", "/admin/maintenance/snmp-trap", datamodel.GetSNMPTrapContent, false)
	e.HTML("GET", "/admin/maintenance/worm-protection", datamodel.GetWormProtectionContent, false)
	e.HTML("GET", "/admin/maintenance/ddos-protection", datamodel.GetDDoSProtectionContent, false)
	e.HTML("GET", "/admin/maintenance/arp-protection", datamodel.GetARPProtectionContent, false)
	e.HTML("GET", "/admin/maintenance/sessions", datamodel.GetSessionsContent, false)

	// 配置模块
	e.HTML("GET", "/admin/config/ports", datamodel.GetPortsContent, false)
	e.HTML("GET", "/admin/config/link-aggregation", datamodel.GetLinkAggregationContent, false)
	e.HTML("GET", "/admin/config/storm-control", datamodel.GetStormControlContent, false)
	e.HTML("GET", "/admin/config/flow-control", datamodel.GetFlowControlContent, false)
	e.HTML("GET", "/admin/config/port-isolation", datamodel.GetPortIsolationContent, false)
	e.HTML("GET", "/admin/config/port-monitor", datamodel.GetPortMonitorContent, false)
	e.HTML("GET", "/admin/config/vlan", datamodel.GetVLANContent, false)
	e.HTML("GET", "/admin/config/mac-table", datamodel.GetMacTableContent, false)
	e.HTML("GET", "/admin/config/stp", datamodel.GetSTPContent, false)
	e.HTML("GET", "/admin/config/erps", datamodel.GetERPSContent, false)
	e.HTML("GET", "/admin/config/poe", datamodel.GetPoEContent, false)
	e.HTML("GET", "/admin/config/port-mirror", datamodel.GetPortMirrorContent, false)
	e.HTML("GET", "/admin/config/multicast", datamodel.GetMulticastContent, false)
	e.HTML("GET", "/admin/config/resource", datamodel.GetResourceContent, false)
	e.HTML("GET", "/admin/config/stack", datamodel.GetStackContent, false)

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
