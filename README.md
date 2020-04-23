## Simple Application written in Golang with MySqllite Database
> Golang 1.14 | Sqllite | EchoFramework | Elastic APM | ZAP Logger | Docker 

To Run the pre-built docker image, that automatically builds with every new push
```
docker run -p 8080:8080 -p 8081:8081 --name xeroApi techievee/xero:v1.0.0 /xeroProductAPI
```
The application exposes 
*  8080 - For standard HTTP port
*  8081 - For TLS port ( TLS Enabled by default from the config with self-signed certificate)

## Configuration
The application configuration can be specified as YAML and their config location can be specified using the -cnf environment variable, which defaults to current directory

- config folder
  - app.yaml
    - app_env - prod: All debug logs are supressed in stdout, any other values: all logs enabled
    - services - For specifying the port and TLS options
  - mysqlite.yaml
    - readwrite-db - Settings for running the write instance connection for mysqlite, Can have only 1 active connection
    - readonly-db  - Settings for running immutable instance connection for mysqlite, Can have only any number of active connection
  - mysqlite_test.yaml
    - All setting to run the the unit testing, similar to mysqlite


## Building the solution

For building the solution, please 
- Install the gcc and other developer tools
- Copy the cert and config folders
- create a folder where the db file resides 
- RUN 
    - CGO_ENABLED=1 go build -o /xeroProductAPI

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

GET Endpoints returns their respective objects

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

Other endpoints, returns the ID of the object
```
 "5fafad6c-ba7f-448a-bd7f-430d986e2e46"
```


