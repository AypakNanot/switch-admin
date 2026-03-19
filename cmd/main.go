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
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/language"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/chartjs"
	"github.com/GoAdminGroup/themes/adminlte"
	"github.com/gin-gonic/gin"
	"switch-admin/internal/datamodel"
	systemDatamodel "switch-admin/internal/datamodel/system"
	maintDatamodel "switch-admin/internal/datamodel/maintenance"
	configDatamodel "switch-admin/internal/datamodel/config"
	networkDatamodel "switch-admin/internal/datamodel/network"
	diagnosticDatamodel "switch-admin/internal/datamodel/diagnostic"
	"switch-admin/internal/handler"
	maintHandler "switch-admin/internal/handler/maintenance"
	networkHandler "switch-admin/internal/handler/network"
	configHandler "switch-admin/internal/handler/config"
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
	maintenanceHandler := maintHandler.New()
	networkHandler := networkHandler.New()
	configHandler := configHandler.New()

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
				"mode":             currentMode,
				"mode_description": modeDesc,
				"database":         "SQLite3 (data/admin.db)",
				"goadmin_version":  "v1.2.26",
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
	r.GET("/api/v1/config/ports", configHandler.GetPorts)
	r.GET("/api/v1/config/ports/:port_id", configHandler.GetPortDetail)
	r.PUT("/api/v1/config/ports/:port_id", configHandler.UpdatePort)

	// 配置模块 API - 链路聚合
	r.GET("/api/v1/link-aggregation", configHandler.GetLinkAggregation)
	r.POST("/api/v1/link-aggregation", configHandler.CreateLinkAggregation)
	r.PUT("/api/v1/link-aggregation/:id", configHandler.UpdateLinkAggregation)
	r.DELETE("/api/v1/link-aggregation/:id", configHandler.DeleteLinkAggregation)

	// 网络模块 API - VLAN 管理
	r.GET("/api/v1/network/vlans", networkHandler.GetVLANList)
	r.POST("/api/v1/network/vlans", networkHandler.CreateVLAN)
	r.PUT("/api/v1/network/vlans/:id", networkHandler.UpdateVLAN)
	r.DELETE("/api/v1/network/vlans/:id", networkHandler.DeleteVLAN)
	r.DELETE("/api/v1/network/vlans", networkHandler.DeleteVLANs)
	r.POST("/api/v1/network/vlans/:id/ports", networkHandler.AddVLANPort)
	r.DELETE("/api/v1/network/vlans/:id/ports", networkHandler.RemoveVLANPort)

	// 网络模块 API - 端口管理
	r.GET("/api/v1/network/ports", networkHandler.GetPortList)
	r.GET("/api/v1/network/ports/:name", networkHandler.GetPortDetail)
	r.PUT("/api/v1/network/ports/:name", networkHandler.UpdatePort)
	r.POST("/api/v1/network/ports/:name/reset", networkHandler.ResetPort)
	r.POST("/api/v1/network/ports/:name/restart", networkHandler.RestartPort)

	// 网络模块 API - 链路聚合管理
	r.GET("/api/v1/network/lags", networkHandler.GetLAGList)
	r.POST("/api/v1/network/lags", networkHandler.CreateLAG)
	r.PUT("/api/v1/network/lags/:id", networkHandler.UpdateLAG)
	r.DELETE("/api/v1/network/lags/:id", networkHandler.DeleteLAG)
	r.POST("/api/v1/network/lags/:id/ports", networkHandler.AddLAGPort)
	r.DELETE("/api/v1/network/lags/:id/ports", networkHandler.RemoveLAGPort)

	// 网络模块 API - STP 管理
	r.GET("/api/v1/network/stp/config", networkHandler.GetSTPConfig)
	r.PUT("/api/v1/network/stp/config", networkHandler.UpdateSTPConfig)
	r.GET("/api/v1/network/stp/status", networkHandler.GetSTPStatus)

	// 网络模块 API - ACL 管理
	r.GET("/api/v1/network/acls", networkHandler.GetACLList)
	r.POST("/api/v1/network/acls", networkHandler.CreateACL)
	r.PUT("/api/v1/network/acls/:id", networkHandler.UpdateACL)
	r.DELETE("/api/v1/network/acls/:id", networkHandler.DeleteACL)
	r.GET("/api/v1/network/acls/:id/rules", networkHandler.GetACLRules)
	r.POST("/api/v1/network/acls/:id/rules", networkHandler.AddACLRule)
	r.PUT("/api/v1/network/acls/:id/rules/:ruleID", networkHandler.UpdateACLRule)
	r.DELETE("/api/v1/network/acls/:id/rules/:ruleID", networkHandler.DeleteACLRule)

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
		UrlPrefix: "admin",
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
	systemDatamodel.InitDatabaseTables(e.SqliteConnection())

	if err := e.AddGenerators(systemDatamodel.Generators).
		AddGenerator("user", systemDatamodel.GetUserTable).
		AddDisplayFilterXssJsFilter().
		Use(r); err != nil {
		panic(err)
	}

	// 初始化菜单
	systemDatamodel.InitMenu(e.SqliteConnection())
	systemDatamodel.InitDashboard(e.SqliteConnection())
	systemDatamodel.InitConfigMenu(e.SqliteConnection())

	r.Static("/uploads", "./uploads")

	// 注册自定义页面 - 使用 GoAdmin 的 HTML 方法
	// Dashboard 和系统配置
	e.HTML("GET", "/admin/dashboard", systemDatamodel.GetDashboardContent, false)
	e.HTML("GET", "/admin/system/config", systemDatamodel.GetSystemConfigPage, false)

	// 网络模块 - IP 路由
	e.HTML("GET", "/admin/network/route-table", networkDatamodel.GetRouteTableContent, false)
	e.HTML("GET", "/admin/network/static-route", networkDatamodel.GetStaticRouteContent, false)

	// 网络模块 - 诊断工具
	e.HTML("GET", "/admin/network/ping", diagnosticDatamodel.GetPingContent, false)
	e.HTML("GET", "/admin/network/traceroute", diagnosticDatamodel.GetTracerouteContent, false)
	e.HTML("GET", "/admin/network/cable-test", diagnosticDatamodel.GetCableTestContent, false)

	// 维护模块
	e.HTML("GET", "/admin/maintenance/reboot-save", maintDatamodel.GetRebootSaveContent, false)
	e.HTML("GET", "/admin/maintenance/users", maintDatamodel.GetUsersContent, false)
	e.HTML("GET", "/admin/maintenance/system-config", maintDatamodel.GetMaintenanceSystemConfigContent, false)
	e.HTML("GET", "/admin/maintenance/load-config", maintDatamodel.GetLoadConfigContent, false)
	e.HTML("GET", "/admin/maintenance/files", maintDatamodel.GetFilesContent, false)
	e.HTML("GET", "/admin/maintenance/logs", maintDatamodel.GetLogsContent, false)
	e.HTML("GET", "/admin/maintenance/snmp", maintDatamodel.GetSNMPContent, false)
	e.HTML("GET", "/admin/maintenance/snmp-trap", maintDatamodel.GetSNMPTrapContent, false)
	e.HTML("GET", "/admin/maintenance/worm-protection", maintDatamodel.GetWormProtectionContent, false)
	e.HTML("GET", "/admin/maintenance/ddos-protection", maintDatamodel.GetDDoSProtectionContent, false)
	e.HTML("GET", "/admin/maintenance/arp-protection", maintDatamodel.GetARPProtectionContent, false)
	e.HTML("GET", "/admin/maintenance/sessions", maintDatamodel.GetSessionsContent, false)

	// 网络模块
	e.HTML("GET", "/admin/network/vlan", networkDatamodel.GetVLANContent, false)
	e.HTML("GET", "/admin/network/port", networkDatamodel.GetPortContent, false)
	e.HTML("GET", "/admin/network/lag", networkDatamodel.GetLAGContent, false)
	e.HTML("GET", "/admin/network/stp", networkDatamodel.GetSTPContent, false)
	e.HTML("GET", "/admin/network/acl", networkDatamodel.GetACLContent, false)

	// 配置模块
	e.HTML("GET", "/admin/config/ports", configDatamodel.GetPortsContent, false)
	e.HTML("GET", "/admin/config/link-aggregation", configDatamodel.GetLinkAggregationContent, false)
	e.HTML("GET", "/admin/config/storm-control", configDatamodel.GetStormControlContent, false)
	e.HTML("GET", "/admin/config/flow-control", configDatamodel.GetFlowControlContent, false)
	e.HTML("GET", "/admin/config/port-isolation", configDatamodel.GetPortIsolationContent, false)
	e.HTML("GET", "/admin/config/port-monitor", configDatamodel.GetPortMonitorContent, false)
	e.HTML("GET", "/admin/config/vlan", configDatamodel.GetVLANContent, false)
	e.HTML("GET", "/admin/config/mac-table", configDatamodel.GetMacTableContent, false)
	e.HTML("GET", "/admin/config/stp", configDatamodel.GetSTPContent, false)
	e.HTML("GET", "/admin/config/erps", configDatamodel.GetERPSContent, false)
	e.HTML("GET", "/admin/config/poe", configDatamodel.GetPoEContent, false)
	e.HTML("GET", "/admin/config/port-mirror", configDatamodel.GetPortMirrorContent, false)
	e.HTML("GET", "/admin/config/multicast", configDatamodel.GetMulticastContent, false)
	e.HTML("GET", "/admin/config/resource", configDatamodel.GetResourceContent, false)
	e.HTML("GET", "/admin/config/stack", configDatamodel.GetStackContent, false)

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
