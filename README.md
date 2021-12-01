# RideSharing (ETI-Assignment1)

## Folder Structure

|       Codebase       |                Description                |
| :------------------: | :---------------------------------------: |
| [frontend](frontend) |          React Next.js Frontend           |
|  [backend](backend)  | Microservices connected to MySQL Database |
|   [server](server)   |    General-Purpose API Backend Server     |

## Usage

There is a total of **6 servers** need to be started. You should be able to see `"Listening on port x ..."` after firing up the servers.

### Microservices Backend

Refer to this [README guide](backend/README.md) to start the 3 microsevice servers + 1 express server which acts as anti-corruption layer.

### General-Purpose API Backend

Refer to this [README guide](server/README.md) to start the express servers which acts as general-purpose API backend to redirect client requests to respective microservices.

### Frontend

Refer to this [README guide](frontend/README.md) to start the next.js frontend.

## Architecture Diagram

![Architecture Diagram](docs/architecture_diagram.png)

## Backend Architecture Design Consideration

### Microservices design

The business logic from the assignment case studies breaks down into **three microservices** (driver, passenger, trip).

In general, each microservice is connected to the MySQL database. The microservice will serve a server endpoint which allows external http requests with methods (GET/PUT/POST) to retrieve, create, or update the row(s) in the backend database related to the microservice module. For example, **driver microservice** server allows for retrieving all filtered drivers with customisable query parameters at `api/v1/drivers`, and retrieving, creating, updating for a specific driver at `api/v1/driver/{driverid}`.

Each of the microservice is **fully independent from one and another**. The three microservices **serve at different localhost endpoint with different port numbers**, where the driver, passenger and trip microservices serves at 8080, 8081 and 8082 port numbers respectively.

In the case when one microservice needs to communicate with another microservices, **anti-corruption layer** is implemented between the layer in the backend.

### Anti-corruption layer
