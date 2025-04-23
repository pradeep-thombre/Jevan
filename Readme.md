# Jevan - Cart, Product, and Order Management System

## Overview

This is a Go-based application that provides APIs for managing carts, products, and orders in a mess management system. The application uses MongoDB as the database and provides Swagger documentation for easy exploration of the APIs.

## Run Locally

Clone the project

```bash
  git clone https://github.com/pradeep-thombre/Jevan.git
```

Go to the project directory

```bash
  cd Jevan
```

Import dependencies

```bash
  go mod tidy
```

Start the server

```bash
  go run .
```

## API Reference

### Cart APIs

#### Add Item to Cart

```http
  POST /cart/id
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Cart ID            |

Payload:
```json
{
    "item_id": "string",       // required
    "quantity": integer,       // required
    "price": float,            // required
    "name": "string"           // required
}
```
Add an item to the specified cart.

#### Get Cart by ID

```http
  GET /cart/:id
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Cart ID           |

Get all items in the specified cart.

#### Delete Item from Cart

```http
  DELETE /cart/id
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Cart ID           |
| `itemId`  | `string` | **Required**. Item ID to delete |

Delete a specific item from the cart.

#### Delete All Items from Cart

```http
  DELETE /cart/id/all
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Cart ID           |

Delete all items from the cart.

### Product APIs

#### Create Product

```http
  POST /products
```

Payload:
```json
{
    "name": "string",        // required
    "description": "string", // required
    "price": float,          // required
    "quantity": integer      // required
}
```
Create a new product and return the product ID.

#### Get All Products

```http
  GET /products
```

Get a list of all products in the store.

#### Get Product by ID

```http
  GET /products/:id
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Product ID         |

Get the details of a specific product by its ID.

#### Update Product by ID

```http
  PUT /products/:id
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Product ID         |

Payload:
```json
{
    "name": "string",        // required
    "description": "string", // required
    "price": float,          // required
    "quantity": integer      // required
}
```

Update the product details by provided ID and payload.

#### Delete Product by ID

```http
  DELETE /products/:id
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Product ID         |

Delete the product by the specified ID.

### Order APIs

#### Create Order

```http
  POST /orders
```

Payload:
```json
{
    "cart_id": "string",      // required
    "total_amount": float,    // required
    "status": "string"        // required (e.g., pending, completed)
}
```
Create a new order using the cart ID and total amount.

#### Get Order by ID

```http
  GET /orders/:id
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Order ID           |

Get the details of a specific order.

#### Update Order Status

```http
  PUT /orders/:id
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Order ID           |

Payload:
```json
{
    "status": "string"  // required (e.g., shipped, cancelled)
}
```

Update the status of an existing order.

#### Cancel Order

```http
  DELETE /orders/:id
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Order ID           |

Cancel an order by ID.

## Swagger Documentation

Swagger UI: [http://localhost:3000/swagger/index.html](http://localhost:3000/swagger/index.html)

## Authors

- [Pradeep Thombre](https://www.github.com/pradeep-thombre)

## 🛠 Tech Stacks

- Golang
- MongoDB
- Swagger APIs
- Echo Framework

## Support

For support, email us.

- [Pradeep Thombre](mailto:padiee755@gmail.com)
