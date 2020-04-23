package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"

	"github.com/techievee/xero/database"
	productServiceCmds "github.com/techievee/xero/productService/commands"
	productServiceCtl "github.com/techievee/xero/productService/controller"
	"github.com/techievee/xero/productService/models"
	"github.com/techievee/xero/xeroHelper"
	"github.com/techievee/xero/xeroLog/debugcore"
)

var pCmd *productServiceCmds.ProductsCmds
var pCtl *productServiceCtl.ProductsCtl
var uuid, po_uuid, p2_uuid string
var testFile string

func TestMain(m *testing.M) {

	// Init Config
	var (
		config     *viper.Viper
		configPath = "./config"
		err        error
	)

	xeroHelper.ParseFlags(configPath)

	config, err = xeroHelper.LoadConfig()
	if err != nil || config == nil {
		log.Errorf("Failed to load config")
		os.Exit(1)
	}

	// Init DB
	db := database.NewDB(config, "mysqlite_test", &debugcore.NoOpsLogger{})
	if db == nil {
		log.Errorf("Failed to load DB")
		os.Exit(1)
	}
	testFile = config.GetString("mysqlite_test.readwrite-db.filepath") + config.GetString("mysqlite_test.readwrite-db.database") + ".db"

	// Init Product Cmd
	pCmd = &productServiceCmds.ProductsCmds{
		DB:     db,
		Logger: &debugcore.NoOpsLogger{},
	}

	// Flush all the data
	dbRw := pCmd.DB.RW(context.Background())
	dbRw.Exec("DELETE FROM ProductOptions")
	dbRw.Exec("DELETE FROM Products")

	pCtl = &productServiceCtl.ProductsCtl{
		ServiceCommands: pCmd,
		Logger:          &debugcore.NoOpsLogger{},
	}

	// Create one product for testing
	p1 := models.Product{
		Name:          "Name P1",
		Description:   "Description P1",
		Price:         10.5,
		DeliveryPrice: 1.5,
	}
	uuid, _ = pCmd.AddNewProduct(context.Background(), p1)

	// Create one product for testing delete
	p2 := models.Product{
		Name:          "Name P2",
		Description:   "Description P2",
		Price:         10.5,
		DeliveryPrice: 1.5,
	}
	p2_uuid, _ = pCmd.AddNewProduct(context.Background(), p2)

	// Create one product option for testing
	// Add new 2 products- options
	po1 := models.ProductOption{
		Name:        "color",
		Description: "Black",
	}
	po_uuid, _ = pCmd.AddNewProductOption(context.Background(), uuid, po1)

	c := m.Run()
	println(testFile)
	os.Remove(testFile)
	os.Exit(c)
}

// Test Database commands

func TestCommands(t *testing.T) {

	ctx := context.Background()

	// Clear the table contents

	// Test Add new products
	p1 := models.Product{
		Name:          "test name",
		Description:   "test description",
		Price:         10.5,
		DeliveryPrice: 1.5,
	}
	uuid, err := pCmd.AddNewProduct(ctx, p1)
	if err != nil {
		t.Error(err)
		return
	}
	// Asset as valid ID
	if len(uuid) == 36 {
		t.Log("Newly inserted record:", uuid)
	} else {
		t.Errorf("Invalid UUID: %v", uuid)
		return
	}

	// Test Update Product
	p2 := models.Product{
		Name:          "test updated",
		Description:   "test description updated",
		Price:         104.5,
		DeliveryPrice: 10.5,
	}
	updateCount, err := pCmd.UpdateProduct(ctx, p2, uuid)
	if err != nil {
		t.Error(err)
	}
	// Asset as valid ID
	if updateCount == 01 {
		t.Log("Updated")
	} else {
		t.Errorf("Not Updated")
	}

	// Test Show all commands
	prod, err := pCmd.FetchAllProducts(ctx, "", "")
	if err != nil {
		t.Error(err)
	}
	if len(prod) == 1 {
		t.Logf("Fetched %d records", len(prod))
	} else {
		t.Errorf("Wrong number of records %d", len(prod))
	}

	// Test with the Updated product name
	prodU, err := pCmd.FetchAllProducts(ctx, "updated", "")
	if err != nil {
		t.Error(err)
	}

	if len(prodU) == 1 {
		t.Logf("Fetched %d records", len(prodU))
	} else {
		t.Errorf("Wrong number of records %d", len(prodU))
	}

	// Test with the old product name
	prodO, err := pCmd.FetchAllProducts(ctx, "name", "")
	if err != nil {
		t.Error(err)
	}
	if len(prodO) == 0 {
		t.Logf("Fetched %d records", len(prodO))
	} else {
		t.Errorf("Wrong number of records %d", len(prodO))
	}

	// Test with the old product name
	prodI, err := pCmd.FetchAllProducts(ctx, "", uuid)
	if err != nil {
		t.Error(err)
	}
	if len(prodI) == 1 {
		t.Logf("Fetched %d records", len(prodI))
	} else {
		t.Errorf("Wrong number of records %d", len(prodI))
	}

	// Add new 2 products- options
	po1 := models.ProductOption{
		Name:        "color",
		Description: "Black",
	}
	po_uuid, err := pCmd.AddNewProductOption(ctx, uuid, po1)
	if err != nil {
		t.Error(err)
		return
	}
	if len(po_uuid) == 36 {
		t.Log("Newly inserted product option:", po_uuid)
	} else {
		t.Errorf("Invalid UUID: %v", po_uuid)
		return
	}

	po2 := models.ProductOption{
		Name:        "color",
		Description: "White",
	}
	po2_uuid, err := pCmd.AddNewProductOption(ctx, uuid, po2)
	if err != nil {
		t.Error(err)
		return
	}
	if len(po_uuid) == 36 {
		t.Log("Newly inserted product option:", po2_uuid)
	} else {
		t.Errorf("Invalid UUID: %v", po2_uuid)
		return
	}

	// Fetch all product option
	prodO1, err := pCmd.FetchAllProductOptions(ctx, uuid, "")
	if err != nil {
		t.Error(err)
	}
	if len(prodO1) == 2 {
		t.Logf("Fetched %d product options", len(prodO1))
	} else {
		t.Errorf("Wrong number of records %d", len(prodO1))
	}

	// Fetch based on product id
	prodO2, err := pCmd.FetchAllProductOptions(ctx, uuid, po_uuid)
	if err != nil {
		t.Error(err)
	}
	if len(prodO2) == 1 {
		t.Logf("Fetched %d product options", len(prodO2))
	} else {
		t.Errorf("Wrong number of records %d", len(prodO2))
	}

	// Update the product option
	po3 := models.ProductOption{
		Name:        "color",
		Description: "Black-Updated",
	}
	pou_count, err := pCmd.UpdateProductOption(ctx, uuid, po_uuid, po3)
	if err != nil {
		t.Error(err)
		return
	}
	if pou_count == 1 {
		t.Log("Deleted ")
	} else {
		t.Errorf("Not deleted the product Option")
		return
	}

	// Delete one option
	poCount, err := pCmd.DeleteProductOption(ctx, uuid, po_uuid)
	if err != nil {
		t.Error(err)
		return
	}
	if poCount == 1 {
		t.Log("Deleted ")
	} else {
		t.Errorf("Not deleted the product Option")
		return
	}

	// Fetch the product option count
	prodO2, err = pCmd.FetchAllProductOptions(ctx, uuid, "")
	if err != nil {
		t.Error(err)
	}
	if len(prodO2) == 1 {
		t.Logf("Fetched %d product options", len(prodO2))
	} else {
		t.Errorf("Wrong number of records %d", len(prodO2))
	}

	// Delete the product
	poDel, err := pCmd.DeleteProduct(ctx, uuid)
	if err != nil {
		t.Error(err)
		return
	}
	if poDel == 1 {
		t.Log("Product Deleted ")
	} else {
		t.Errorf("Not deleted the product")
		return
	}

	// Check the total count
	prod, err = pCmd.FetchAllProducts(ctx, "", "")
	if err != nil {
		t.Error(err)
	}
	if len(prod) == 0 {
		t.Logf("Fetched %d records", len(prod))
	} else {
		t.Errorf("Wrong number of records %d", len(prod))
	}

}

// Test Product controllers

func TestAddNewProduct(t *testing.T) {

	productJson :=
		`
		{
		  "Name": "iPhone SE",
		  "Description": "Updated Second Gen Version.",
		  "Price": 1229.99,
		  "DeliveryPrice": 1.99
		}
		`

	e := echo.New()
	request := httptest.NewRequest(http.MethodPost, "/api/products", strings.NewReader(productJson))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	responseRecorder := httptest.NewRecorder()
	c := e.NewContext(request, responseRecorder)
	pCtl.AddNewProduct(c)
	if responseRecorder.Code != http.StatusCreated {
		t.Logf("Expected : %d\n got:%d\n", http.StatusCreated, responseRecorder.Code)
		t.Fail()
	}
	body := responseRecorder.Body.String()
	t.Logf("Output: %v", body)

}

func TestAddNewProductInvalidPrice(t *testing.T) {

	productJson :=
		`
		{
		  "Name": "iPhone SE",
		  "Description": "Updated Second Gen Version.",
		  "Price": -1229.99,
		  "DeliveryPrice": 1.99
		}
		`

	e := echo.New()
	request := httptest.NewRequest(http.MethodGet, "/api/products", strings.NewReader(productJson))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	responseRecorder := httptest.NewRecorder()
	c := e.NewContext(request, responseRecorder)
	pCtl.AddNewProduct(c)
	if responseRecorder.Code != http.StatusBadRequest {
		t.Logf("Expected : %d\n got:%d\n", http.StatusBadRequest, responseRecorder.Code)
		t.Fail()
	}
	body := responseRecorder.Body.String()
	t.Logf("Output: %v", body)

}

func TestShowAllProducts(t *testing.T) {

	e := echo.New()
	request := httptest.NewRequest(http.MethodGet, "/api/products", strings.NewReader(""))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	responseRecorder := httptest.NewRecorder()
	c := e.NewContext(request, responseRecorder)
	pCtl.ShowProducts(c)
	if responseRecorder.Code != http.StatusOK {
		t.Logf("Expected : %d\n got:%d\n", http.StatusOK, responseRecorder.Code)
		t.Fail()
	}
	body := responseRecorder.Body.String()
	t.Logf("Output: %v", body)

}

func TestShowAllProductsWithName(t *testing.T) {

	e := echo.New()
	request := httptest.NewRequest(http.MethodGet, "/api/products?name=samsung", strings.NewReader(""))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	responseRecorder := httptest.NewRecorder()
	c := e.NewContext(request, responseRecorder)
	pCtl.ShowProducts(c)
	if responseRecorder.Code != http.StatusOK {
		t.Logf("Expected : %d\n got:%d\n", http.StatusOK, responseRecorder.Code)
		t.Fail()
	}
	body := responseRecorder.Body.String()
	t.Logf("Output: %v", body)

}

func TestShowProductWithID(t *testing.T) {

	e := echo.New()
	request := httptest.NewRequest(http.MethodPost, "/api/products/", strings.NewReader(""))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	responseRecorder := httptest.NewRecorder()
	c := e.NewContext(request, responseRecorder)
	c.SetPath("/api/products/:id")
	c.SetParamNames("id")
	c.SetParamValues(uuid)
	pCtl.ShowProduct(c)
	if responseRecorder.Code != http.StatusOK {
		t.Logf("Expected : %d\n got:%d\n", http.StatusOK, responseRecorder.Code)
		t.Fail()
	}
	body := responseRecorder.Body.String()
	t.Logf("Output: %v", body)

}

func TestShowAllProductWithName(t *testing.T) {

	e := echo.New()
	q := make(url.Values)
	q.Set("name", "p1")
	request := httptest.NewRequest(http.MethodGet, "/api/products/?"+q.Encode(), strings.NewReader(""))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	responseRecorder := httptest.NewRecorder()
	c := e.NewContext(request, responseRecorder)
	pCtl.ShowProducts(c)
	if responseRecorder.Code != http.StatusOK {
		t.Logf("Expected : %d\n got:%d\n", http.StatusOK, responseRecorder.Code)
		t.Fail()
	}
	body := responseRecorder.Body.String()
	t.Logf("Output: %v", body)

}

func TestUpdateProductsWithID(t *testing.T) {

	productJson :=
		`
		{
		  "Name": "iPhone SE- Updated",
		  "Description": "Updated Second Gen Version.",
		  "Price": 1229.99,
		  "DeliveryPrice": 1.99
		}
		`

	e := echo.New()
	request := httptest.NewRequest(http.MethodPut, "/api/products/", strings.NewReader(productJson))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	responseRecorder := httptest.NewRecorder()
	c := e.NewContext(request, responseRecorder)
	c.SetPath("/api/products/:id")
	c.SetParamNames("id")
	c.SetParamValues(uuid)
	pCtl.UpdateProduct(c)
	if responseRecorder.Code != http.StatusOK {
		t.Logf("Expected : %d\n got:%d\n", http.StatusOK, responseRecorder.Code)
		t.Fail()
	}
	body := responseRecorder.Body.String()
	t.Logf("Output: %v", body)

}

func TestUpdateProductsWithIDInvalidID(t *testing.T) {

	productJson :=
		`
		{
		  "Name": "iPhone SE- Updated",
		  "Description": "Updated Second Gen Version.",
		  "Price": 1229.99,
		  "DeliveryPrice": 1.99
		}
		`

	e := echo.New()
	request := httptest.NewRequest(http.MethodPut, "/api/products/", strings.NewReader(productJson))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	responseRecorder := httptest.NewRecorder()
	c := e.NewContext(request, responseRecorder)
	c.SetPath("/api/products/:id")
	c.SetParamNames("id")
	c.SetParamValues("69d6c863-18e4-4f21-8f46-9cc5128a84c4")
	pCtl.UpdateProduct(c)
	if responseRecorder.Code != http.StatusBadRequest {
		t.Logf("Expected : %d\n got:%d\n", http.StatusBadRequest, responseRecorder.Code)
		t.Fail()
	}
	body := responseRecorder.Body.String()
	t.Logf("Output: %v", body)

}

func TestShowProductWithName(t *testing.T) {

	productJson :=
		`
		{
		  "Name": "iPhone SE- Updated",
		  "Description": "Updated Second Gen Version.",
		  "Price": 1229.99,
		  "DeliveryPrice": 1.99
		}
		`

	e := echo.New()
	request := httptest.NewRequest(http.MethodPut, "/api/products/", strings.NewReader(productJson))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	responseRecorder := httptest.NewRecorder()
	c := e.NewContext(request, responseRecorder)
	c.SetPath("/api/products/:id")
	c.SetParamNames("id")
	c.SetParamValues("69d6c863-18e4-4f21-8f46-9cc5128a84c4")
	pCtl.UpdateProduct(c)
	if responseRecorder.Code != http.StatusBadRequest {
		t.Logf("Expected : %d\n got:%d\n", http.StatusBadRequest, responseRecorder.Code)
		t.Fail()
	}
	body := responseRecorder.Body.String()
	t.Logf("Output: %v", body)

}

func TestDeleteProduct(t *testing.T) {

	e := echo.New()
	request := httptest.NewRequest(http.MethodDelete, "/api/products/", strings.NewReader(""))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	responseRecorder := httptest.NewRecorder()
	c := e.NewContext(request, responseRecorder)
	c.SetPath("/api/products/:id")
	c.SetParamNames("id")
	c.SetParamValues(p2_uuid)
	pCtl.DeleteProduct(c)
	if responseRecorder.Code != http.StatusOK {
		t.Logf("Expected : %d\n got:%d\n", http.StatusOK, responseRecorder.Code)
		t.Fail()
	}
	body := responseRecorder.Body.String()
	t.Logf("Output: %v", body)

}

//-- TEST PRODUCT OPTION

func TestAddNewAndUpdateProductOption(t *testing.T) {

	productOptionJson :=
		`
		{
		  "Name": "Color",
		  "Description": "Black"
		}
		`

	e := echo.New()
	request := httptest.NewRequest(http.MethodPost, "/api/products/:id/options/", strings.NewReader(productOptionJson))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	responseRecorder := httptest.NewRecorder()
	c := e.NewContext(request, responseRecorder)
	c.SetParamNames("id")
	c.SetParamValues(uuid)
	pCtl.AddNewProductOption(c)
	if responseRecorder.Code != http.StatusCreated {
		t.Logf("Expected : %d\n got:%d\n", http.StatusCreated, responseRecorder.Code)
		t.Fail()
	}
	body := responseRecorder.Body.String()
	t.Logf("Output: %v", body)

	productOptionJson =
		`
		{
		  "Name": "Color",
		  "Description": "Blue"
		}
		`

	request = httptest.NewRequest(http.MethodPut, "/api/products/:id/options/:optionId", strings.NewReader(productOptionJson))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	responseRecorder = httptest.NewRecorder()
	c = e.NewContext(request, responseRecorder)
	c.SetParamNames("id", "optionId")

	optionId := strings.Trim(body, "/n")
	optionId = strings.Trim(body, `"`)
	c.SetParamValues(uuid, optionId)
	pCtl.AddNewProductOption(c)
	if responseRecorder.Code != http.StatusCreated {
		t.Logf("Expected : %d\n got:%d\n", http.StatusCreated, responseRecorder.Code)
		t.Fail()
	}
	body = responseRecorder.Body.String()
	t.Logf("Output: %v", body)

}

func TestAddNewProductInvalidOption(t *testing.T) {

	productOptionJson :=
		`
		{
		  "Name": "Color",
		  "Description": "Black"
		}
		`

	e := echo.New()
	request := httptest.NewRequest(http.MethodPost, "/api/products/:id/options/", strings.NewReader(productOptionJson))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	responseRecorder := httptest.NewRecorder()
	c := e.NewContext(request, responseRecorder)
	c.SetParamNames("id")
	c.SetParamValues("deed6cfc-9cd8-41fc-b8c0-038f4c1c79cf")
	pCtl.AddNewProductOption(c)
	if responseRecorder.Code != http.StatusBadRequest {
		t.Logf("Expected : %d\n got:%d\n", http.StatusBadRequest, responseRecorder.Code)
		t.Fail()
	}
	body := responseRecorder.Body.String()
	t.Logf("Output: %v", body)

}

func TestShowAllProductOptions(t *testing.T) {

	e := echo.New()
	request := httptest.NewRequest(http.MethodGet, "/api/products/:id/options", strings.NewReader(""))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	responseRecorder := httptest.NewRecorder()
	c := e.NewContext(request, responseRecorder)
	c.SetParamNames("id")
	c.SetParamValues(uuid)
	pCtl.ShowProductOptions(c)

	if responseRecorder.Code != http.StatusOK {
		t.Logf("Expected : %d\n got:%d\n", http.StatusOK, responseRecorder.Code)
		t.Fail()
	}
	body := responseRecorder.Body.String()
	t.Logf("Output: %v", body)

}

func TestShowProductOptions(t *testing.T) {

	e := echo.New()
	request := httptest.NewRequest(http.MethodGet, "/api/products/:id/options/:optionId", strings.NewReader(""))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	responseRecorder := httptest.NewRecorder()
	c := e.NewContext(request, responseRecorder)
	c.SetParamNames("id", "optionId")
	c.SetParamValues(uuid, po_uuid)
	pCtl.ShowProductOption(c)

	if responseRecorder.Code != http.StatusOK {
		t.Logf("Expected : %d\n got:%d\n", http.StatusOK, responseRecorder.Code)
		t.Fail()
	}
	body := responseRecorder.Body.String()
	t.Logf("Output: %v", body)

}

func TestDeleteProductOption(t *testing.T) {

	e := echo.New()
	request := httptest.NewRequest(http.MethodGet, "/api/products/:id/options/:optionId", strings.NewReader(""))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	responseRecorder := httptest.NewRecorder()
	c := e.NewContext(request, responseRecorder)
	c.SetParamNames("id", "optionId")
	c.SetParamValues(uuid, po_uuid)
	pCtl.DeleteProductOption(c)
	if responseRecorder.Code != http.StatusOK {
		t.Logf("Expected : %d\n got:%d\n", http.StatusOK, responseRecorder.Code)
		t.Fail()
	}
	body := responseRecorder.Body.String()
	t.Logf("Output: %v", body)

}
