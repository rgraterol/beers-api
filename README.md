# Beers-API
![technology Go](https://img.shields.io/badge/technology-go-blue.svg)
![coverage](https://img.shields.io/badge/coverage-95%25-brightgreen)
![rating](https://img.shields.io/badge/rating-★★★★★-brightgreen)
![codacy](https://img.shields.io/badge/codacy-A-brightgreen)

CRUD API for handling beers in a bar ecosystem.

# Technologies
- Go 1.17
- MySQL


# To run the project
 - Initialize the DB
```
docker-compose up -d
```

- Download dependencies
```bash
go mod download
```

- Run Project
```bash
go run cmd/api/main.go
```

- Test if server is up and running
```bash
curl --location --request GET 'http://localhost:8080/ping'
``` 

## Tests

To test the application run the following command

````bash
go test  ./... -covermode=atomic  -coverpkg=./... -count=1  -race -timeout=30m
````

# Endpoints

### Create `POST /beers`
Persist a beer inside a DB. 

Beer object
```go
type Beer struct {
	ID        int64 
	Name      string
	Brewery   string
	Country   string
	Price     float64
	Currency  string
}
```

Inside the API can only be one beer for each name, brewery and country. Example:
```json
{
  "name": "Ambar",
  "brewery": "Cuello Negro",
  "country": "Chile"
}
```
Cannot be repeated. But there can be other from a different brewery like
```json
{
  "name": "Ambar",
  "brewery": "La Cibeles",
  "country": "Chile"
}
```

The fields Name, Price, Currency.

#### cURL Example
```bash
curl --location --request POST 'http://localhost:8080/beers' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name":"Calafate",
    "brewery":"Austral",
    "price": 1023.432,
    "currency": "ARS"
}'
```

### List `GET /beers`
Retrieves all the beers inside the DB inside a lis of beers.

#### cURL Example
```bash
curl --location --request GET 'http://localhost:8080/beers' \
--header 'Content-Type: application/json'
```

Response
```json
[
    {
        "id": 22,
        "name": "Golden",
        "brewery": "",
        "country": "Chile",
        "price": 100.4,
        "currency": "USD"
    },
    {
        "id": 23,
        "name": "Calafate",
        "brewery": "Austral",
        "country": "Chile",
        "price": 1023.432,
        "currency": "ARS"
    }
]
```

### Get `GET /beers/{beerID}`
Retrieves a single .

#### cURL Example
```bash
curl --location --request GET 'http://localhost:8080/beers/22' \
--header 'Content-Type: application/json'
```

Response
```json
{
    "id": 22,
    "name": "Golden",
    "brewery": "",
    "country": "Chile",
    "price": 100.4,
    "currency": "USD"
}
```

### BoxPrice `GET /beers/{beerID}/boxprice?currency=USD&quantity=4`
Retrieves the price of the desired beer specified by the URL param `beerID`
It accepts two optional query params
- Currency
- Quantity (default:6 if value is not specified)
```go
type BeerBoxParameters struct {
	Currency string `json:"currency"`
	Quantity int64  `json:"quantity"`
}
```

Responds an BoxPrice object
```go
type BeerBox struct {
	Price  float64           `json:"price"`
	Target BeerBoxParameters `json:"target"`
	Beer   Beer              `json:"beer"`
}
```
If goes agains the API of [https://currencylayer.com/](https://currencylayer.com/) which gives the current conversion rate between currencies.
Example response (shorthand version):
```json
{
    "success": true,
    "terms": "https://currencylayer.com/terms",
    "privacy": "https://currencylayer.com/privacy",
    "timestamp": 1644135065,
    "source": "USD",
    "quotes": {
        "USDCLP": 828.503912,
        "USDARS": 105.356594,
        "USDALL": 106.703989
    }
}
```

#### cURL Example
```bash
curl --location --request GET 'http://localhost:8080/beers/23/boxprice?currency=CLP' \
--header 'Content-Type: application/json' 
```

Response with the convertion rates of:
- Target currency CLP to Dollar: 828.503912
- Stored currency ARS to Dollar: 105.356594
- Default box amount of 6 beers per box.
- Price of the beer in ARS 1023.432.

`((828.503912) / (105.356594)) * 1023.432 * 6 = 48288.42980626257`

```json
{
  "price": 48288.42980626257,
  "target": {
    "currency": "CLP",
    "quantity": 6
  },
  "beer": {
    "id": 23,
    "name": "Calafate",
    "brewery": "Austral",
    "country": "",
    "price": 1023.432,
    "currency": "ARS"
  }
}
```

