# RideSharing Backend Server

## Description

This folder contains all the backend logic written in microservices, and a server acts as anti-corruption layer to handle and redirect internal microservices requests.

## Folder Structure

|                      Codebase                      |                                        Description                                         |
| :------------------------------------------------: | :----------------------------------------------------------------------------------------: |
|    [Driver Microservice](microservices/drivers)    |                  `Driver` microservice server connected to MySQL database                  |
| [Passenger Microservice](microservices/passengers) |                `Passenger` microservice server connected to MySQL database                 |
|      [Trip Microservice](microservices/trips)      |                   `Trip` microservice server connected to MySQL database                   |
|      [Internal Microservices Server](server)       | A server acts as an `anti-corruption layer` <br> to handle internal microservices requests |
