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
