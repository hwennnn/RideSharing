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

In general, each microservice is connected to the same MySQL database. **(Disclaimer: In real practice, each microservice should be connected to its own MySQL database. This is done to simply fulfill the assignment purposes.)** The microservice will serve a server endpoint which allows external http requests with methods (GET/PUT/POST) to retrieve, create, or update the row(s) in the backend database related to the microservice module. For example, **driver microservice** server allows for retrieving all filtered drivers with customisable query parameters at `api/v1/drivers`, and retrieving, creating, updating for a specific driver at `api/v1/driver/{driverid}`.

Each of the microservice is **fully independent from one and another**. The three microservices **serve at different localhost endpoint with different port numbers**, where the driver, passenger and trip microservices serves at 8080, 8081 and 8082 port numbers respectively.

In the case when one microservice needs to communicate with another microservices, **anti-corruption layer** is implemented between the layer in the backend.

### Anti-corruption layer

In betweeen the microservices, the **anti-corruption layer** is implemented to handle and redirect internal microservices requests. The server is written in **node.js and Express.js**.

By definition, anti-corruption layer (ACL) creates an isolating layer to provide clients with functionality in terms of their own domain model. The layer talks to the other system through its existing interface, requiring little or no modification to the other system. Internally, the layer translates in both directions as necessary between the two models. (Quoted from **Eric Evans, Domain Driven Design, 16th printing, page 365**)

In short, an ACL is not just about separating the messy code, but as a means to communicate between bounded contexts. It translates from one context to the other, so that data in each context reflects the language and the way that that context thinks and talks about the data

In the project, the express server serves at `localhost:4000`. When one microservice wants to communicate with another microservice, the microservice only needs to send http request with related pathname, for example `api/v1/drivers`, to the server which acts as ACL. By this way the microservice would not have to explicitly call that particular microservice endpoint, which is against the domain driven design (DDD), but instead, only knowing and calling one surface endpoint which will do the lifting communication work for the origin microservice.

The usage of ACL in the project is that when the trip microservice needs to communicate with both driver and passenger microservice in order to retrieve the driver and passenger information for the trip. The trip microservice will send http request with related information to the ACL, and the ACL will redirect the requests to respective microservices and return the desired result back to the trip microservice.

In addition, the ACL will also **authenticate the bearer token** to ensure the requests are valid and sent from the legitimate microservices server. When no or incorrect token is sent, that http request will be blocked, and 403 status code which indicates access forbidden will be sent back.
