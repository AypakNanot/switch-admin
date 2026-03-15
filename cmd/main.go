package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"time"

	_ "github.com/GoAdminGroup/go-admin/adapter/gin"
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/sqlite"
	_ "github.com/GoAdminGroup/themes/sword"

	"github.com/GoAdminGroup/go-admin/engine"
	"switch-admin/internal/datamodel"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/language"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/chartjs"
	"github.com/GoAdminGroup/themes/adminlte"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

	r := gin.New()

	e := engine.Default()

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

	if err := e.AddConfig(&cfg).
		AddGenerators(datamodel.Generators).
		AddGenerator("user", datamodel.GetUserTable).
		AddDisplayFilterXssJsFilter().
		Use(r); err != nil {
		panic(err)
	}

	r.Static("/uploads", "./uploads")

	// 自定义页面 - Dashboard
	e.HTML("GET", "/admin", datamodel.GetDashboardContent)

	go func() {
		_ = r.Run(":9033")
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Print("closing database connection")
	e.SqliteConnection().Close()
}
