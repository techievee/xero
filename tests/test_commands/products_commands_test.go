package test_commands

import (
	"context"
	"log"
	"os"
	"testing"

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
		configPath = "./../config"
		err        error
	)

	xeroHelper.ParseFlags(configPath)

	config, err = xeroHelper.LoadConfig()
	if err != nil || config == nil {
		log.Fatal("Failed to load config")
		os.Exit(1)
	}

	// Init DB
	db := database.NewDB(config, "mysqlite_test", &debugcore.NoOpsLogger{})
	if db == nil {
		log.Fatal("Failed to load DB")
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

	c := m.Run()
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
