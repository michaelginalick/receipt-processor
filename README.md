# receipt-processor

## How to run the application

- To run the application, ensure Docker is installed on your machine, and then follow these steps:
  - Build the Docker image:
    - docker build -t receipt-processor-app .
  - Run the Docker container, exposing port 8080:
    - docker run -p 8080:8080 receipt-processor-app

Once running, you can interact with the application via curl. For example:

To process a receipt:
```bash
curl -d '{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },{
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },{
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },{
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },{
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}' -H "Content-Type: application/json" -X POST http://localhost:8080/receipts/process
```
Make note of the ID returned, which you can use to calculate points:
```bash
curl localhost:8080/receipts/<ID>/points
```

## Run the tests
- Find the running container ID:
    - ```docker ps```
  - Exec into the container
    - ```docker exec -it <container id> /bin/sh ```
    - ``` go test -v ./... ```

## Testing

For testing, the standard Go testing library is utilized to keep dependencies minimal. However, I personally favor [ginkgo](https://onsi.github.io/ginkgo/) for its readability in spec-style tests. Clear and understandable tests are crucial for maintaining project health.

## In memory DB

The persistence layer uses a simple `sync.Map` wrapped in an interface with two methods. While this approach may be over engineered, it offers a readable abstraction that enhances code clarity in my opinion.

## Receipt

The `Receipt` package drives the application's logic, employing a rule set represented by function pointers. This modular approach isolates rule changes, aiming to simplify maintenance and extensibility.
