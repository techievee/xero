## Simple Application written in Golang with MySqllite Database

The application uses Labstack/Echo framework.

It connects to the backend MySqllite Instance and performs a select query.

## API Endpoints

| SNo |           ENDPOINTS                | REST API | METHOD |                       DESCRIPTION                             |
|-----|------------------------------------|----------|--------|---------------------------------------------------------------|
|  1  | /products                          | Yes      |  GET   | gets all products.                                            |
|  2  | /products?name={name}              | Yes      |  GET   | finds all products matching the specified name.               |
|  3  | /products/{:id}                    | Yes      |  GET   | gets the project that matches the specified ID - ID is a GUID.|
|  4  | /products                          | Yes      |  POST  | creates a new product.                                        |
|  5  | /products/{:id}                    | Yes      |  PUT   | updates a product.                                            |
|  6  | /products/{:id}                    | Yes      |  DELETE| deletes a product and its options.                            |
|  7  | /products/{id}/options             | Yes      |  GET   | finds all options for a specified product.                    |
|  8  | /products/{:id}/options/{:optionId}| Yes      |  GET   | finds the specified product option for the specified product. |
|  9  | /products/{:id}/options            | Yes      |  POST  | adds a new product option to the specified product.           |
| 10  | /products/{:id}/options/{:optionId}| Yes      |  PUT   | updates the specified product option.                         |
| 11  | /products/{:id}/options/{:optionId}| Yes      |  DELETE| deletes the specified product option.                         |



## Data Models

**Product:**
```
{
  "Id": "01234567-89ab-cdef-0123-456789abcdef",
  "Name": "Product name",
  "Description": "Product description",
  "Price": 123.45,
  "DeliveryPrice": 12.34
}
```

**Products:**
```
{
  "Items": [
    {
      // product
    },
    {
      // product
    }
  ]
}
```

**Product Option:**
```
{
  "Id": "01234567-89ab-cdef-0123-456789abcdef",
  "Name": "Product name",
  "Description": "Product description"
}
```

**Product Options:**
```
{
  "Items": [
    {
      // product option
    },
    {
      // product option
    }
  ]
}
```





