package commands

import (
	"context"
	"database/sql"
	"strings"

	"github.com/google/uuid"
	"go.elastic.co/apm"

	"github.com/techievee/xero/productService/models"
)

const (
	stmtProducts      = "SELECT Id, Name, Description, Price, DeliveryPrice FROM Products"
	stmtInsertProduct = "INSERT INTO  Products (Id, Name, Description, Price, DeliveryPrice) VALUES (?,?,?,?,?)"
	stmtUpdateProduct = "UPDATE Products SET Name=?, Description=?, Price=?, DeliveryPrice=? WHERE Id=? COLLATE NOCASE"
	stmtDeleteProduct = "DELETE FROM Products WHERE Id=? COLLATE NOCASE"
)

func (c *ProductsCmds) FetchAllProducts(ctx context.Context, pName string, pID string) ([]models.DBProducts, error) {

	span, ctx := apm.StartSpan(ctx, "products.show", "db")
	span.SpanData.Context.SetTag("span", "FetchAllProducts")
	defer span.End()

	db := c.DB.RO(ctx)

	params := ""
	stmt := stmtProducts
	if pID != "" {
		stmt += " WHERE Id=? COLLATE NOCASE"
		params = strings.ToLower(pID)
	} else if pName != "" {
		stmt += " WHERE Name like ? COLLATE NOCASE "
		params = "%" + strings.ToLower(pName) + "%"
	}

	var rows *sql.Rows
	var err error
	if params != "" {
		rows, err = db.QueryContext(ctx, stmt, params)
	} else {
		rows, err = db.QueryContext(ctx, stmt)
	}
	if err != nil {
		c.Logger.Error("Error while fetching products", "error", err)
		return nil, err
	}

	result := []models.DBProducts{}
	for rows.Next() {
		dbObj := models.DBProducts{}
		rows.Scan(&dbObj.DBID, &dbObj.DBName, &dbObj.DBDescription, &dbObj.DBPrice, &dbObj.DBDeliveryPrice)
		result = append(result, dbObj)
	}
	if err = rows.Err(); err != nil {
		c.Logger.Error("Error while scanning rows", "error", err)
		return nil, err
	}

	c.Logger.Debug("Fetched all the products", "total_rows", len(result))
	return result, nil
}

func (c *ProductsCmds) AddNewProduct(ctx context.Context, product models.Product) (string, error) {

	span, ctx := apm.StartSpan(ctx, "products.add", "db")
	span.SpanData.Context.SetTag("span", "AddNewProduct")
	defer span.End()

	db := c.DB.RW(ctx)
	id := uuid.New()
	statement, _ := db.Prepare(stmtInsertProduct)
	result, err := statement.ExecContext(ctx, id, product.Name, product.Description, product.Price, product.DeliveryPrice)
	if err != nil {
		c.Logger.Error("Error while inserting new rows", "error", err)
		return "", err
	}
	insertedID, _ := result.LastInsertId()
	c.Logger.Debug("Added new product", "inserted_id", insertedID, "uuid", id)
	return id.String(), err
}

func (c *ProductsCmds) UpdateProduct(ctx context.Context, product models.Product, productID string) (int64, error) {

	span, ctx := apm.StartSpan(ctx, "products.update", "db")
	span.SpanData.Context.SetTag("span", "UpdateProduct")
	defer span.End()

	db := c.DB.RW(ctx)
	statement, _ := db.Prepare(stmtUpdateProduct)
	result, err := statement.ExecContext(ctx, product.Name, product.Description, product.Price, product.DeliveryPrice, productID)
	if err != nil {
		c.Logger.Error("Error while updating products", "error", err)
		return 0, err
	}
	affectedRows, err := result.RowsAffected()
	c.Logger.Debug("Updated the products", "affected_rows", affectedRows)
	return affectedRows, err
}

func (c *ProductsCmds) DeleteProduct(ctx context.Context, productID string) (int64, error) {

	span, ctx := apm.StartSpan(ctx, "products.delete", "db")
	span.SpanData.Context.SetTag("span", "DeleteProduct")
	defer span.End()

	db := c.DB.RW(ctx)

	// Delete all the options related to this product
	if _, err := c.DeleteAllProductOptions(ctx, productID); err != nil {
		c.Logger.Error("Error while deleting product options", "error", err)
		return 0, err
	}

	statement, _ := db.Prepare(stmtDeleteProduct)
	result, err := statement.ExecContext(ctx, productID)
	if err != nil {
		c.Logger.Error("Error while deleting products", "error", err)
		return 0, err
	}
	affectedRows, _ := result.RowsAffected()
	c.Logger.Debug("Deleted the product", "affected_rows", affectedRows)
	return affectedRows, err
}
