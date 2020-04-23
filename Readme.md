## Simple Application written in Golang with MySqllite Database

The application uses Labstack/Echo framework.

It connects to the backend MySqllite Instance and performs a select query.

## API Endpoints

| SNo |           ENDPOINTS                | REST API | METHOD |                       DESCRIPTION                             |
|-----|------------------------------------|----------|--------|---------------------------------------------------------------|
|  1  | /products                          | Yes      |  GET   | gets all products.                                            |
|  2  | /products?name={name}              | Yes      |  GET   | finds all products matching the specified name.               |
|  3  | /products/{:id}                    | Yes      |  GET   | gets the product that matches the specified ID - ID GUID/UUID.|
|  4  | /products                          | Yes      |  POST  | creates a new product.                                        |
|  5  | /products/{:id}                    | Yes      |  PUT   | updates the product with specified ID.                        |
|  6  | /products/{:id}                    | Yes      |  DELETE| deletes a product and its options.                            |
|  7  | /products/{id}/options             | Yes      |  GET   | finds all options for a specified product.                    |
|  8  | /products/{:id}/options/{:optionId}| Yes      |  GET   | finds the specified product option for the specified product. |
|  9  | /products/{:id}/options            | Yes      |  POST  | adds a new product option to the specified product.           |
| 10  | /products/{:id}/options/{:optionId}| Yes      |  PUT   | updates the specified product option.                         |
| 11  | /products/{:id}/options/{:optionId}| Yes      |  DELETE| deletes the specified product option.                         |


## API Return Code

| SNo |           ENDPOINTS                | REST API | METHOD |                       DESCRIPTION                             |
|-----|------------------------------------|----------|--------|---------------------------------------------------------------|
|  1  | /products                          | Yes      |  GET   | 200- Success, 500- Internal Server Error.                     |
|  2  | /products?name={name}              | Yes      |  GET   | 200- Success, 500- Internal Server Error.                     |
|  3  | /products/{:id}                    | Yes      |  GET   | 200- Success, 500- Internal Server Error, 400- Invalid ID     |
|  4  | /products                          | Yes      |  POST  | 201- Successfully created, 500- Server Err, 400- Invalid data |
|  5  | /products/{:id}                    | Yes      |  PUT   | 200- Success, 500- Internal Server Error, 400- Invalid ID     |
|  6  | /products/{:id}                    | Yes      |  DELETE| 200- Success, 500- Internal Server Error, 400- Invalid ID     |
|  7  | /products/{id}/options             | Yes      |  GET   | 200- Success, 500- Internal Server Error, 400- Invalid ID     |
|  8  | /products/{:id}/options/{:optionId}| Yes      |  GET   | 200- Success, 500- Internal Server Error, 400- Invalid ID     |
|  9  | /products/{:id}/options            | Yes      |  POST  | 201- Successfully created, 500- Server Err, 400- Invalid data |
| 10  | /products/{:id}/options/{:optionId}| Yes      |  PUT   | 200- Success, 500- Internal Server Error, 400- Invalid ID     |
| 11  | /products/{:id}/options/{:optionId}| Yes      |  DELETE| 200- Success, 500- Internal Server Error, 400- Invalid ID     |


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





