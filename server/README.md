# RideSharing General-Purpose API Backend Server

## Description

The server is written in **node.js and Express.js**. The server acts as **general-purpose API backend server** which **redirects the client requests to respective microservice**. With this server, there will be **only one server endpoint surfaced** for the clients because the server will do the heavy work by **redirecting the traffic to the correct microservice** where **each of them is served at different ports**. Besides, the server will also **authenticate the token** to ensure the requests are valid and sent from the micorservices server.

### Prerequisite

You need npm installed

```
npm install npm@latest -g
```

### Install all the packages

```bash
npm install
npm ci
```

### Usage

```bash
npm run dev
```
