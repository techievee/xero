package ctls

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"go.elastic.co/apm"

	"github.com/techievee/xero/productService/models"
	xError "github.com/techievee/xero/xeroErrors"
	"github.com/techievee/xero/xeroHelper"
)

func (p *ProductsCtl) ShowProductOptions(c echo.Context) error {

	defer xError.CatchErr(nil)
	ctx := c.Request().Context()
	span, _ := apm.StartSpan(ctx, "products_options.show", "api")
	defer span.End()

	// Look for the id param
	productId := c.Param("id")
	if strings.Trim(productId, " ") == "" || xeroHelper.ValidateUUID(productId) == false {
		// Return 400, Bad request
		return c.JSON(http.StatusBadRequest, "Invalid Request : Valid Product Id required")
	}

	// Check for existence of Product
	if product, err := p.ServiceCommands.FetchAllProducts(ctx, "", productId); err != nil {
		return xError.NewUnexpectedGenericError(err)
	} else if len(product) == 0 {
		return c.JSON(http.StatusBadRequest, "Invalid Request : Valid Product Id required")
	}

	items := []models.ProductOption{}
	result, err := p.ServiceCommands.FetchAllProductOptions(ctx, productId, "")
	if err != nil {
		if err != sql.ErrNoRows {
			// Return 500, Server Error
			return xError.NewUnexpectedGenericError(err)
		}
	}

	if len(result) > 0 {
		for _, v := range result {
			item := models.ProductOption{}

			// Safely convert the DbTypes to GoTypes
			if v.DBID.Valid {
				item.ID = v.DBID.String
			}
			if v.DBName.Valid {
				item.Name = v.DBName.String
			}
			if v.DBDescription.Valid {
				item.Description = v.DBDescription.String
			}
			items = append(items, item)

		}

	} else {
		items = []models.ProductOption{}
	}

	resultProductOptions := models.ProductOptions{
		Items: &items,
	}

	// Return 200
	return c.JSON(http.StatusOK, resultProductOptions)

}

func (p *ProductsCtl) ShowProductOption(c echo.Context) error {

	defer xError.CatchErr(nil)
	ctx := c.Request().Context()
	span, _ := apm.StartSpan(ctx, "products_options.show", "api")
	defer span.End()

	// Look for the id param
	productId := c.Param("id")
	if strings.Trim(productId, " ") == "" || xeroHelper.ValidateUUID(productId) == false {
		// Return 400, Bad request
		return c.JSON(http.StatusBadRequest, "Invalid Request : Valid Product Id required")
	}

	// Check for existence of Product
	if product, err := p.ServiceCommands.FetchAllProducts(ctx, "", productId); err != nil {
		return xError.NewUnexpectedGenericError(err)
	} else if len(product) == 0 {
		return c.JSON(http.StatusBadRequest, "Invalid Request : Valid Product Id required")
	}

	productOptionId := c.Param("optionId")
	if xeroHelper.ValidateUUID(productOptionId) == false {
		// Return 400, Bad request
		return c.JSON(http.StatusBadRequest, "Invalid Request Format")
	}

	items := []models.ProductOption{}
	result, err := p.ServiceCommands.FetchAllProductOptions(ctx, productId, productOptionId)
	if err != nil {
		if err != sql.ErrNoRows {
			// Return 500, Server Error
			return xError.NewUnexpectedGenericError(err)
		}
	}

	if len(result) == 0 {
		// Return 400, Bad request
		return c.JSON(http.StatusBadRequest, "Invalid option Id")
	}
	productOption := models.ProductOption{}
	for _, v := range result {

		// Safely convert the DbTypes to GoTypes
		if v.DBID.Valid {
			productOption.ID = v.DBID.String
		}
		if v.DBName.Valid {
			productOption.Name = v.DBName.String
		}
		if v.DBDescription.Valid {
			productOption.Description = v.DBDescription.String
		}
		items = append(items, productOption)

	}

	// Return 200
	return c.JSON(http.StatusOK, productOption)

}

func (p *ProductsCtl) AddNewProductOption(c echo.Context) error {

	defer xError.CatchErr(nil)
	ctx := c.Request().Context()
	span, _ := apm.StartSpan(ctx, "product_option.add", "api")
	defer span.End()

	productId := c.Param("id")
	if strings.Trim(productId, " ") == "" || xeroHelper.ValidateUUID(productId) == false {
		// Return 400, Bad request
		return c.JSON(http.StatusBadRequest, "Invalid Request : Valid Product Id required")
	}

	// Check for existence of Product
	if product, err := p.ServiceCommands.FetchAllProducts(ctx, "", productId); err != nil {
		return xError.NewUnexpectedGenericError(err)
	} else if len(product) == 0 {
		return c.JSON(http.StatusBadRequest, "Invalid Request : Valid Product Id required")
	}

	// Parse the productOption from the post body
	productOption := models.ProductOption{}
	if err := c.Bind(&productOption); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Request Format")
	}

	// Validate the format of the productOption json
	if err := productOption.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Request Format :"+err.Error())
	}

	// Validate the name
	id, err := p.ServiceCommands.AddNewProductOption(ctx, productId, productOption)
	if err != nil {
		return xError.NewUnexpectedGenericError(err)
	}

	// Return 200 with Newly created ID
	return c.JSON(http.StatusCreated, id)

}

func (p *ProductsCtl) UpdateProductOption(c echo.Context) error {

	defer xError.CatchErr(nil)
	ctx := c.Request().Context()
	span, _ := apm.StartSpan(ctx, "product_option.update", "api")
	defer span.End()

	productId := c.Param("id")
	if !xeroHelper.ValidateUUID(productId) {
		return c.JSON(http.StatusBadRequest, "Invalid product id")
	}
	// Check for existence of Product
	if product, err := p.ServiceCommands.FetchAllProducts(ctx, "", productId); err != nil {
		return xError.NewUnexpectedGenericError(err)
	} else if len(product) == 0 {
		return c.JSON(http.StatusBadRequest, "Invalid Request : Valid Product Id required")
	}

	productOptionId := c.Param("optionId")
	if xeroHelper.ValidateUUID(productOptionId) == false {
		// Return 400, Bad request
		return c.JSON(http.StatusBadRequest, "Invalid Request Format")
	}

	// Parse the productOption of the post parameter
	productOption := models.ProductOption{}
	if err := c.Bind(&productOption); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Request Format")
	}

	// Validate the format of the productOption json
	if err := productOption.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Request Format"+err.Error())
	}

	// Validate the name
	affectedRows, err := p.ServiceCommands.UpdateProductOption(ctx, productId, productOptionId, productOption)
	if err != nil {
		// Returns 500, Server error
		return xError.NewUnexpectedGenericError(err)
	}
	if affectedRows == 0 {
		return c.JSON(http.StatusBadRequest, "Invalid product Option id")
	}

	return c.JSON(http.StatusOK, productOptionId)

}

func (p *ProductsCtl) DeleteProductOption(c echo.Context) error {

	defer xError.CatchErr(nil)
	ctx := c.Request().Context()
	span, _ := apm.StartSpan(ctx, "product_option.delete", "api")
	defer span.End()

	productId := c.Param("id")
	if !xeroHelper.ValidateUUID(productId) {
		return c.JSON(http.StatusBadRequest, "Invalid productOption id")
	}
	// Check for existence of Product
	if product, err := p.ServiceCommands.FetchAllProducts(ctx, "", productId); err != nil {
		return xError.NewUnexpectedGenericError(err)
	} else if len(product) == 0 {
		return c.JSON(http.StatusBadRequest, "Invalid Request : Valid Product Id required")
	}

	productOptionId := c.Param("optionId")
	if xeroHelper.ValidateUUID(productOptionId) == false {
		// Return 400, Bad request
		return c.JSON(http.StatusBadRequest, "Invalid Request Format")
	}

	// Validate the name
	affectedRows, err := p.ServiceCommands.DeleteProductOption(ctx, productId, productOptionId)
	if err != nil {
		// Returns 500, Server error
		return xError.NewUnexpectedGenericError(err)
	}
	if affectedRows == 0 {
		return c.JSON(http.StatusBadRequest, "Invalid product option id")
	}

	return c.JSON(http.StatusOK, productOptionId)

}
