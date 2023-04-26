# E-commerce Backend Logic with Golang

This repository contains the backend logic for an e-commerce application built using the Go programming language.

### Instalation
To install this project, clone the repository:
```bash
git clone https://github.com/Xenn-00/go-merce.git
```
Then, install the required dependencies:
```bash
go mod download
```
### Usage
To run the project, use the following command:
```bash
go run main.go
```
This will start the server on `localhost:7000`.

### API Endpoints
The following API endpoints are available:
| Method | Endpoint |
|----- | ------ |
| POST | localhost:7000/users/signup |
| POST | localhost:7000/users/login |
| POST | localhost:7000/admin/addProduct |
| GET | localhost:7000/users/productView |
| GET | locahost:7000/users/search | 
| GET | localhost:7000/addtocart |
| GET | localhost:7000/removeitem |
| GET | localhost:7000/listcart |
| POST | localhost:7000/addaddress |
| PUT | localhost:7000/edithomeaddress |
| PUT | localhost:7000/editworkaddress |
| DELETE | localhost:7000/deleteaddresses |
| GET | localhost:7000/cartcheckout |
| GET | localhost:7000/instantbuy |
