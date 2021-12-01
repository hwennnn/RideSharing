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

### Prerequisite

You would need MySQL installed. Refer to this [guide](https://dev.mysql.com/doc/mysql-installation-excerpt/5.7/en/) on how to download and install MySQL on your machine.

### Execute MySQL setup script

The script helps to create the database and tables. It also helps to setup some initial dummy data for the demo purposes. If u **encounter errors when executing the script**, its likely u already have a database called RideSharing, you will need to drop that table before executing the command below.

```bash
mysql> source path/to/setup.sql
```

### Start microservice server

Refer to this [README guide](database/README.md) to start the 3 microsevices servers.

### Start express server

Refer to this [README guide](server/README.md) to start the express server which acts as anti-corruption layer
