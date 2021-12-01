# RideSharing Backend Server

## Description

This folder contains all the backend logic written in microservices, and a server acts as anti-corruption layer to handle and redirect internal microservices requests.

## Folder Structure

|                          Codebase                           |                                        Description                                         |
| :---------------------------------------------------------: | :----------------------------------------------------------------------------------------: |
|    [Driver Microservice](database/microservices/drivers)    |                  `Driver` microservice server connected to MySQL database                  |
| [Passenger Microservice](database/microservices/passengers) |                `Passenger` microservice server connected to MySQL database                 |
|      [Trip Microservice](database/microservices/trips)      |                   `Trip` microservice server connected to MySQL database                   |
|           [Internal Microservices Server](server)           | A server acts as an `anti-corruption layer` <br> to handle internal microservices requests |

## Usage

### Start microservice server

Refer to this [README guide](database/README.md) to start the 3 microsevice servers.

### Start express server

Refer to this [README guide](server/README.md) to start the express server which acts as anti-corruption layer
