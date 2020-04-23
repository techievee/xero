package ctls

import (
	"net/http"

	"github.com/labstack/echo"
	"go.elastic.co/apm"

	"github.com/techievee/xero/productService/models"
	xError "github.com/techievee/xero/xeroErrors"
	"github.com/techievee/xero/xeroHelper"
)

func (p *ProductsCtl) ShowProducts(c echo.Context) error {

	defer xError.CatchErr(nil)
	ctx := c.Request().Context()
	span, _ := apm.StartSpan(ctx, "products.show", "api")
	defer span.End()

	// Look for the name param
	productName := c.QueryParam("name")

	items := []models.Product{}
	result, err := p.ServiceCommands.FetchAllProducts(ctx, productName, "")
	if err != nil {
		return xError.NewUnexpectedGenericError(err)
	}

	if len(result) > 0 {
		for _, v := range result {
			item := models.Product{}

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
			if v.DBPrice.Valid {
				item.Price = v.DBPrice.Float64
			}
			if v.DBDeliveryPrice.Valid {
				item.DeliveryPrice = v.DBDeliveryPrice.Float64
			}
			items = append(items, item)

		}

	} else {
		items = []models.Product{}
	}

	resultProducts := models.Products{
		Items: &items,
	}

	// Return 200
	return c.JSON(http.StatusOK, resultProducts)

}

func (p *ProductsCtl) ShowProduct(c echo.Context) error {

	defer xError.CatchErr(nil)
	ctx := c.Request().Context()
	span, _ := apm.StartSpan(ctx, "product.show", "api")
	defer span.End()

	// Look for the id param
	productId := c.Param("id")
	if xeroHelper.ValidateUUID(productId) == false {
		// Return 400, Bad request
		return c.JSON(http.StatusBadRequest, "Invalid Request Format")
	}

	result, err := p.ServiceCommands.FetchAllProducts(ctx, "", productId)
	if err != nil {
		return xError.NewUnexpectedGenericError(err)
	}
	if len(result) == 0 {
		return c.JSON(http.StatusBadRequest, "Product ID not available")
	}

	product := models.Product{}
	for _, v := range result {

		// Safely convert the DbTypes to GoTypes
		if v.DBID.Valid {
			product.ID = v.DBID.String
		}
		if v.DBName.Valid {
			product.Name = v.DBName.String
		}
		if v.DBDescription.Valid {
			product.Description = v.DBDescription.String
		}
		if v.DBPrice.Valid {
			product.Price = v.DBPrice.Float64
		}
		if v.DBDeliveryPrice.Valid {
			product.DeliveryPrice = v.DBDeliveryPrice.Float64
		}

	}

	// Return 200
	return c.JSON(http.StatusOK, product)

}

func (p *ProductsCtl) AddNewProduct(c echo.Context) error {

	defer xError.CatchErr(nil)
	ctx := c.Request().Context()
	span, _ := apm.StartSpan(ctx, "product.add", "api")
	defer span.End()

	// Parse the product from the post body
	product := models.Product{}
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Request Format")
	}

	// Validate the format of the product json
	if err := product.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Request Format"+err.Error())
	}

	// Validate the name
	id, err := p.ServiceCommands.AddNewProduct(ctx, product)
	if err != nil {
		return xError.NewUnexpectedGenericError(err)
	}

	// Return 200 with Newly created ID
	return c.JSON(http.StatusCreated, id)

}

func (p *ProductsCtl) UpdateProduct(c echo.Context) error {

	defer xError.CatchErr(nil)
	ctx := c.Request().Context()
	span, _ := apm.StartSpan(ctx, "product.update", "api")
	defer span.End()

	productId := c.Param("id")
	if !xeroHelper.ValidateUUID(productId) {
		return c.JSON(http.StatusBadRequest, "Invalid product id")
	}

	// Parse the product of the post parameter
	product := models.Product{}
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Request Format")
	}

	// Validate the format of the product json
	if err := product.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Request Format"+err.Error())
	}

	// Validate the name
	affectedRows, err := p.ServiceCommands.UpdateProduct(ctx, product, productId)
	if err != nil {
		// Returns 500, Server error
		return xError.NewUnexpectedGenericError(err)
	}
	if affectedRows == 0 {
		return c.JSON(http.StatusBadRequest, "Invalid product id")
	}

	return c.JSON(http.StatusOK, productId)

}

func (p *ProductsCtl) DeleteProduct(c echo.Context) error {

	defer xError.CatchErr(nil)
	ctx := c.Request().Context()
	span, _ := apm.StartSpan(ctx, "product.delete", "api")
	defer span.End()

	productId := c.Param("id")
	if !xeroHelper.ValidateUUID(productId) {
		return c.JSON(http.StatusBadRequest, "Invalid product id")
	}

	// Validate the name
	affectedRows, err := p.ServiceCommands.DeleteProduct(ctx, productId)
	if err != nil {
		// Returns 500, Server error
		return xError.NewUnexpectedGenericError(err)
	}
	if affectedRows == 0 {
		return c.JSON(http.StatusBadRequest, "Invalid product id")
	}

	return c.JSON(http.StatusOK, productId)

}
