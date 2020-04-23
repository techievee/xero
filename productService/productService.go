package productService

import (
	"github.com/spf13/viper"

	"github.com/techievee/xero/apiServer"
	"github.com/techievee/xero/database"
	productServiceCmds "github.com/techievee/xero/productService/commands"
	productServiceCtl "github.com/techievee/xero/productService/controller"
	"github.com/techievee/xero/xeroLog/debugcore"
)

type ProductService struct {
	Config            *viper.Viper
	ServiceController *productServiceCtl.ProductsCtl

	DB      *database.DB
	RestAPI *apiServer.APIServer
	Logger  debugcore.Logger
}

func NewProductService(config *viper.Viper, db *database.DB, restAPI *apiServer.APIServer, logger debugcore.Logger) *ProductService {

	// Create a new controller for the Product
	productsCmds := &productServiceCmds.ProductsCmds{DB: db, Logger: logger}
	productsCtl := &productServiceCtl.ProductsCtl{ServiceCommands: productsCmds, Logger: logger}

	return &ProductService{
		Config:            config,
		ServiceController: productsCtl,
		RestAPI:           restAPI,
		Logger:            logger,
	}
}

func (ps *ProductService) SetupService() {

	ps.Logger.Debug("Product Service Starting")
	ps.LoadRoutes()

}

func (ps *ProductService) LoadRoutes() {

	ps.Logger.Debug("Setting up routes")
	productsRoute := ps.RestAPI.EchoFramework.Group("/api/products")

	// Products Routes
	productsRoute.GET("", ps.ServiceController.ShowProducts)
	productsRoute.GET("/:id", ps.ServiceController.ShowProduct)
	productsRoute.POST("", ps.ServiceController.AddNewProduct)
	productsRoute.PUT("/:id", ps.ServiceController.UpdateProduct)
	productsRoute.DELETE("/:id", ps.ServiceController.DeleteProduct)

	// ProductOption Routes
	productsRoute.GET("/:id/options", ps.ServiceController.ShowProductOptions)
	productsRoute.GET("/:id/options/:optionId", ps.ServiceController.ShowProductOption)
	productsRoute.POST("/:id/options", ps.ServiceController.AddNewProductOption)
	productsRoute.PUT("/:id/options/:optionId", ps.ServiceController.UpdateProductOption)
	productsRoute.DELETE("/:id/options/:optionId", ps.ServiceController.DeleteProductOption)

	ps.Logger.Debug("Routes were successfully configured")
}
