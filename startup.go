package main

import (
	"os"

	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/spf13/viper"

	"github.com/techievee/xero/apiServer"
	"github.com/techievee/xero/database"
	"github.com/techievee/xero/productService"
	"github.com/techievee/xero/xeroHelper"
	"github.com/techievee/xero/xeroLog"
	"github.com/techievee/xero/xeroLog/debugcore"
)

var AppFramework *echo.Echo

// Global Variable used by all the other packages

func main() {

	var (
		config     *viper.Viper
		configPath = "D:\\Xero\\xero\\config"
		err        error
	)

	xeroHelper.ParseFlags(configPath)
	glog.Infof("Starting Product API")

	// Configuration Initialization
	// Parse flag should have loaded the config, if its not loaded then load it via the loadconfig functions
	// Load the configuration files
	if config == nil {
		if config, err = xeroHelper.LoadConfig(); err != nil || config == nil {
			glog.Error("Failed to load config") // REVIEW : Should this panic
			return
		}
	}

	// Logger Initialization
	// Logger in injected to all the service, to maintain the logs
	env := config.GetString("app.app_env")
	xeroLogger := xeroLog.NewLogger(env, xeroLog.WithServiceName("xero-api"))
	xeroLogger.Debug("xeroLogger successfuly configured")

	// Database initialization
	xeroLogger.Debug("Initializing the DB")
	db := database.NewDB(config, "mysqlite", xeroLogger)
	if db == nil {
		xeroLogger.Debug("Error initializing DB")
		os.Exit(1)
	}
	xeroLogger.Debug("Successfully initialized DB")

	// Starting the API framework for serving the Prodcut
	xeroLogger.Debug("Initializing the Rest Framework")
	restAPI := apiServer.NewRestAPI(env, config, xeroLogger)

	xeroLogger.Debug("Starting Products API Service")
	startProductsService(config, db, restAPI, xeroLogger)

	go restAPI.StartServer()

	tlsEnabled := config.Get("app.service.tls.enabled").(bool)
	if tlsEnabled {
		go restAPI.StartTLSServer()
	}

	quit := make(chan bool)
	<-quit

}

func startProductsService(config *viper.Viper, db *database.DB, restAPI *apiServer.APIServer, logger debugcore.Logger) {
	ps := productService.NewProductService(config, db, restAPI, logger)
	ps.SetupService()
}
