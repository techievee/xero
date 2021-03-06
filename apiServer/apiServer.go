package apiServer

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
	"go.elastic.co/apm/module/apmecho"

	xError "github.com/techievee/xero/xeroErrors"
	"github.com/techievee/xero/xeroLog/debugcore"
)

const (
	environmentProd = "prod"
)

var Debug bool

type APIServer struct {
	EchoFramework *echo.Echo
	errorHandler  func(err error, c echo.Context)
	Logger        debugcore.Logger
	appConfig     *viper.Viper
}

func NewRestAPI(env string, appConfig *viper.Viper, logger debugcore.Logger) *APIServer {
	echoFramework := echo.New()
	Debug = env != environmentProd
	echoFramework.Debug = Debug
	echoFramework.HideBanner = true

	// Remove trailing slash from the request
	echoFramework.Pre(middleware.RemoveTrailingSlash())

	// Inject a Random requestID to track every request
	echoFramework.Use(middleware.RequestID())

	// Use the APM (Application performance monitoring)
	echoFramework.Use(apmecho.Middleware())

	return &APIServer{
		EchoFramework: echoFramework,
		errorHandler:  HTTPErrorHandler,
		Logger:        logger,
		appConfig:     appConfig,
	}

}

func (s *APIServer) StartServer() {

	server := s.appConfig.GetString("app.service.host") + ":" + s.appConfig.GetString("app.service.port")

	err := s.EchoFramework.Start(server)
	if err != nil {
		s.Logger.Error("Cannot start the Server", "error", err)
	}

}

func (s *APIServer) StartTLSServer() {

	tlsServer := s.appConfig.GetString("app.service.tls.host") + ":" + s.appConfig.GetString("app.service.tls.port")

	err := s.EchoFramework.StartTLS(tlsServer, s.appConfig.GetString("app.service.tls.certificate"), s.appConfig.GetString("app.service.tls.key"))
	if err != nil {
		s.Logger.Error("Cannot start the TLS Server", "error", err)
	}
}

// HTTPErrorHandler handles the error response and sends a valid response to the frontend
// If the Debug is set to prod, then the traceback value is not sent to front-end
func HTTPErrorHandler(err error, c echo.Context) {

	var e xError.Error
	code := http.StatusBadRequest

	switch v := err.(type) {
	case xError.Error:
		e = v
		xError.LogStdError(e)

		// Also log context timeout from the timeout middleware
		if e.Message == "context deadline exceeded" {
			xError.LogStdError(e)
		}

	case *echo.HTTPError:
		e = xError.New(v.Code, v.Error(), xError.Failed)
	case error:
		e = xError.New(code, v.Error(), xError.Failed)

	}

	if !Debug {
		e.Traceback = nil //traceback
	}

	c.JSON(e.Code, e)
}
