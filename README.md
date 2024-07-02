# go-fastest-api-response

This project is a challenge to fetch the fastest response between two different APIs using multithreading in Go. The two requests are made simultaneously, and the result of the fastest request is displayed on the command line.

## Description

The project makes simultaneous requests to the following APIs:
- [BrasilAPI](https://brasilapi.com.br/api/cep/v1/01153000)
- [ViaCEP](http://viacep.com.br/ws/01153000/json/)

The fastest response is displayed on the command line with the address details and which API sent it. If the response time exceeds 1 second, a timeout error is displayed.

## How to Use

### Prerequisites

- [Go](https://golang.org/dl/) installed.

### Installation

Clone the repository:

```sh
git clone https://github.com/pansani/go-fastest-api-response.git
cd go-fastest-api-response
```

Initialize the Go module:

```sh
go mod init github.com/pansani/go-fastest-api-response
```

### Running

Run the project:

```sh
go run main.go
```

## Project Structure

```go
go-fastest-api-response/
├── go.mod
└── main.go
```

### `main.go`

This file contains the main logic of the project. It makes simultaneous requests and displays the fastest result or a timeout error.

## Author

- João Mendonça - [pansani](https://github.com/pansani)
