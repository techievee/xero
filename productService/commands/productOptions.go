package commands

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"go.elastic.co/apm"

	"github.com/techievee/xero/productService/models"
)

const (
	stmtProductOptions         = "SELECT Id, Name, Description FROM ProductOptions WHERE ProductId=? COLLATE NOCASE "
	stmtInsertProductOption    = "INSERT INTO  ProductOptions (Id, ProductId, Name, Description) VALUES (?,?,?,?)"
	stmtUpdateProductOption    = "UPDATE ProductOptions SET Name=?, Description=? WHERE Id=? COLLATE NOCASE and ProductId=? COLLATE NOCASE"
	stmtDeleteProductOption    = "DELETE FROM ProductOptions WHERE Id=? COLLATE NOCASE and ProductId=? COLLATE NOCASE"
	stmtDeleteAllProductOption = "DELETE FROM ProductOptions WHERE ProductId=? COLLATE NOCASE"
)

// Returns all the product option for the specified product id
func (c *ProductsCmds) FetchAllProductOptions(ctx context.Context, pID string, pOptionID string) ([]models.DBProductOptions, error) {

	span, ctx := apm.StartSpan(ctx, "product_options.show", "db")
	span.SpanData.Context.SetTag("span", "ShowProductOptions")
	defer span.End()

	db := c.DB.RO(ctx)

	params := []interface{}{pID}
	stmt := stmtProductOptions
	if pOptionID != "" {
		stmt += " AND Id=? COLLATE NOCASE "
		params = append(params, pOptionID)
	}

	var rows *sql.Rows
	var err error
	rows, err = db.QueryContext(ctx, stmt, params...)
	if err != nil {
		c.Logger.Error("Error while fetching product options", "error", err)
		return nil, err
	}

	result := []models.DBProductOptions{}
	for rows.Next() {
		dbObj := models.DBProductOptions{}
		rows.Scan(&dbObj.DBID, &dbObj.DBName, &dbObj.DBDescription)
		result = append(result, dbObj)
	}
	if err = rows.Err(); err != nil {
		c.Logger.Error("Error while scanning rows", "error", err)
		return nil, err
	}

	c.Logger.Debug("Fetched all the product options", "total_rows", len(result))
	return result, nil
}

// Returns the newly added product option id
func (c *ProductsCmds) AddNewProductOption(ctx context.Context, pID string, product models.ProductOption) (string, error) {

	span, ctx := apm.StartSpan(ctx, "product_options.add", "db")
	span.SpanData.Context.SetTag("span", "AddNewProductOption")
	defer span.End()

	db := c.DB.RW(ctx)
	id := uuid.New()
	statement, _ := db.Prepare(stmtInsertProductOption)
	result, err := statement.ExecContext(ctx, id, pID, product.Name, product.Description)
	if err != nil {
		c.Logger.Error("Error while inserting new rows to product option", "error", err)
		return "", err
	}
	insertedID, _ := result.LastInsertId()
	c.Logger.Debug("Added new product option", "inserted_id", insertedID, "uuid", id)
	return id.String(), err
}

// Returns total number of rows affected by this update
func (c *ProductsCmds) UpdateProductOption(ctx context.Context, pID string, pOptionID string, product models.ProductOption) (int64, error) {

	span, ctx := apm.StartSpan(ctx, "product_options.update", "db")
	span.SpanData.Context.SetTag("span", "UpdateProductOptions")
	defer span.End()

	db := c.DB.RW(ctx)
	statement, _ := db.Prepare(stmtUpdateProductOption)
	result, err := statement.ExecContext(ctx, product.Name, product.Description, pOptionID, pID)
	if err != nil {
		c.Logger.Error("Error while updating product options", "error", err)
		return 0, err
	}
	affectedRows, err := result.RowsAffected()
	c.Logger.Debug("Updated the product options", "affected_rows", affectedRows)
	return affectedRows, err
}

// Delete the product specified in the product option
func (c *ProductsCmds) DeleteProductOption(ctx context.Context, pID string, pOptionID string) (int64, error) {

	span, ctx := apm.StartSpan(ctx, "product_options.delete", "db")
	span.SpanData.Context.SetTag("span", "DeleteProductoption")
	defer span.End()

	db := c.DB.RW(ctx)
	statement, _ := db.Prepare(stmtDeleteProductOption)
	result, err := statement.ExecContext(ctx, pOptionID, pID)
	if err != nil {
		c.Logger.Error("Error while deleting product option", "error", err)
		return 0, err
	}
	affectedRows, _ := result.RowsAffected()
	c.Logger.Debug("Deleted the product option", "affected_rows", affectedRows)
	return affectedRows, err
}

// Delete the all options for the specified product
func (c *ProductsCmds) DeleteAllProductOptions(ctx context.Context, pID string) (int64, error) {

	span, ctx := apm.StartSpan(ctx, "product_options.delete", "db")
	span.SpanData.Context.SetTag("span", "DeleteAllProductoptions")
	defer span.End()

	db := c.DB.RW(ctx)
	statement, _ := db.Prepare(stmtDeleteAllProductOption)
	result, err := statement.ExecContext(ctx, pID)
	if err != nil {
		c.Logger.Error("Error while deleting all product option", "error", err)
		return 0, err
	}
	affectedRows, _ := result.RowsAffected()
	c.Logger.Debug("Deleted all the product options", "affected_rows", affectedRows)
	return affectedRows, err
}
