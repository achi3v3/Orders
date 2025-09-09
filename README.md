# Orders

Demo service with Kafka, PostgreSQL, and cache. This project implements a microservice architecture for processing and storing order data.

## Project Structure

```
.
├── configs/
│   └── .env                    # Environment configuration file
├── front/
│   ├── Dockerfile             # Docker configuration for frontend service
│   ├── index.html            # Main HTML page
│   ├── script.js             # Client-side JavaScript for order lookup
│   └── style.css             # CSS styling
├── main-service/              # Main service for consuming Kafka messages and serving HTTP requests
│   ├── Dockerfile             # Docker configuration for main service
│   ├── go.mod                 # Go module definition
│   ├── go.sum                 # Go dependencies checksum
│   ├── main.go               # Main application entry point
│   ├── internal/
│   │   ├── config/           # Configuration loading utilities
│   │   │   ├── kafka.go     # Kafka configuration
│   │   │   ├── postgres.go # PostgreSQL configuration
│   │   │   └── utils.go     # Environment utilities
│   │   ├── database/        # Database handling
│   │   │   ├── queries.go   # Database queries
│   │   │   └── tables.go    # Database table creation
│   │   ├── models/          # Data models
│   │   │   └── models.go    # Order, Delivery, Payment, and Item structs
│   │   ├── subs/            # Order processing logic
│   │   │   ├── handlers.go  # HTTP handlers
│   │   │   ├── queries.go   # Database queries for orders
│   │   │   ├── repository.go # Data access layer
│   │   │   └── service.go   # Business logic with caching
│   │   ├── kafka/
│   │   │   └── messaging/   # Kafka consumer implementation
│   │   │       ├── consumer.go # Kafka message consumer
│   │   │       └── types.go    # Kafka message types
│   │   └── router/
│   │       └── router.go    # HTTP router setup
└── producer-service/         # Service for producing Kafka messages
    ├── Dockerfile           # Docker configuration for producer service
    ├── go.mod              # Go module definition
    ├── go.sum             # Go dependencies checksum
    ├── main.go           # Main application entry point
    └── internal/
        ├── config/       # Configuration loading utilities
        │   ├── kafka.go # Kafka configuration
        │   ├── postgres.go # PostgreSQL configuration (possibly unused in producer)
        │   └── utils.go    # Environment utilities
        ├── kafka/
        │   └── messaging/  # Kafka producer implementation
        │       ├── producer.go # Kafka message producer
        │       └── types.go    # Kafka message types
        ├── kafka/producer/  # Order message producer
        │   ├── sender.go    # Message sending logic
        │   └── test.go      # Test data generation
        └── models/
            └── models.go    # Data models for order messages
```

## Overview

This project demonstrates a microservice architecture with the following components:

1. **Main Service**:

   - Consumes order messages from Kafka
   - Stores order data in PostgreSQL database
   - Provides HTTP endpoints to retrieve order information
   - Implements in-memory caching for improved performance (cache speeds up the request by ~35 times)
2. **Frontend Service**:

   - Provides a web interface for looking up order information
   - Communicates with the main service via HTTP requests
3. **Producer Service**:

   - Generates and sends order messages to Kafka
4. **Kafka**:

   - Message broker for communication between services
5. **PostgreSQL**:

   - Database for persistent storage of order data

## Services Details

### Main Service

The main service is responsible for:

- Consuming messages from Kafka topic
- Processing and storing order data in PostgreSQL
- Managing an in-memory cache with TTL (Time To Live)
- Providing REST API endpoints for order retrieval

Key components:

- `consumer.go`: Kafka consumer implementation that reads messages and passes them to handlers
- `service.go`: Business logic layer with caching functionality
- `repository.go`: Data access layer for database operations
- `router.go`: HTTP router that exposes endpoints for order retrieval

### Frontend Service

The frontend service provides a simple web interface for users to look up order information:

- Users can enter an order ID in the input field
- The service makes HTTP requests to the main service to retrieve order data
- Results are displayed in a formatted manner on the web page

Key components:

- `index.html`: Main page with input field and result display area
- `script.js`: Client-side JavaScript that handles user input and API requests
- `style.css`: CSS styling for the web interface

### Producer Service

The producer service:

- Generates test order data
- Sends messages to Kafka topic

Key components:

- `sender.go`: Logic for sending messages to Kafka
- `test.go`: Test data generation and sending

## Database Schema

The PostgreSQL database contains the following tables:

- `orders`: Main order information
- `deliveries`: Delivery details for each order
- `payments`: Payment information
- `items`: Items in each order

### Example Data JSON-Scheme

```
{
   "order_uid": "b563feb7b2b84b6test",
   "track_number": "WBILMTESTTRACK",
   "entry": "WBIL",
   "delivery": {
      "name": "Test Testov",
      "phone": "+9720000000",
      "zip": "2639809",
      "city": "Kiryat Mozkin",
      "address": "Ploshad Mira 15",
      "region": "Kraiot",
      "email": "test@gmail.com"
   },
   "payment": {
      "transaction": "b563feb7b2b84b6test",
      "request_id": "",
      "currency": "USD",
      "provider": "wbpay",
      "amount": 1817,
      "payment_dt": 1637907727,
      "bank": "alpha",
      "delivery_cost": 1500,
      "goods_total": 317,
      "custom_fee": 0
   },
   "items": [
      {
         "chrt_id": 9934930,
         "track_number": "WBILMTESTTRACK",
         "price": 453,
         "rid": "ab4219087a764ae0btest",
         "name": "Mascaras",
         "sale": 30,
         "size": "0",
         "total_price": 317,
         "nm_id": 2389212,
         "brand": "Vivienne Sabo",
         "status": 202
      }
   ],
   "locale": "en",
   "internal_signature": "",
   "customer_id": "test",
   "delivery_service": "meest",
   "shardkey": "9",
   "sm_id": 99,
   "date_created": "2021-11-26T06:22:19Z",
   "oof_shard": "1"
}
```

## Configuration

Environment configuration is stored in `configs/.env`:

- PostgreSQL connection details
- Kafka broker URL
- Topic name and consumer group
- Logger level

## Docker Compose

The `docker-compose.yml` file defines all services:

- Kafka service with proper configuration
- PostgreSQL service
- Main service with dependencies
- Frontend service with dependencies
- Producer service with dependencies

All services are containerized and can be run with a single command.

To run the project:

1. Build all services without using cache:
   ```bash
   docker compose build --no-cache
   ```

2. Start all services:
   ```bash
   docker compose up
   ```

3. To run the producer service separately (after the other services are up):
   ```bash
   docker compose run producer-service
   ```

The frontend service will be available at http://localhost:3000, and the main service API at http://localhost:8080.

<img width="1659" height="929" alt="image" src="https://github.com/user-attachments/assets/b258bce7-22af-4bf5-974f-c993e361a2c9" />
