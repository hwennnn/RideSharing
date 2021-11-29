# RideSharing Internal Microservices Server

## Description

The server is written in **node.js and Express.js**. The server acts as **anti-corruption layer** to handle and redirect internal microservices requests. The server will also **authenticate the token** to ensure the requests are valid and sent from the micorservices server.

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
