# RideSharing Microservice API Documentation

**Base URL: http://localhost:80/server/**

## Token Authentication

The server will **authenticate the bearer token** to ensure the requests are valid and sent from the frontend server. When no or incorrect token is sent, that http request will be blocked, and 403 status code which indicates access forbidden will be sent back.

Hence, it is required to supply the bearer token in the authorization header as following:

```js
const requestConfig = {
  headers: { Authorization: `Bearer ${authenticationToken}` },
};
```

## 1. Driver

### Data Structures

|      Field Name       |  Type  |                                         Description                                          |
| :-------------------: | :----: | :------------------------------------------------------------------------------------------: |
|       driver_id       | string |                            The unique ID which identifies driver                             |
|      first_name       | string |                                 The first name of the driver                                 |
|       last_name       | string |                                 The last name of the driver                                  |
|     mobile_number     | string |                               The mobile number of the driver                                |
|     email_address     | string |                               The email address of the driver                                |
| identification_number | string |     The identification number of the driver. It cannot be edited after driver creation.      |
|  car_license_number   | string |                            The car license number of the driver.                             |
|   available_status    |  int   |                          The available status number of the driver.                          |
|                       |        | **0 -> used by golang to indicate whether the integer variable has been initialised or not** |
|                       |        |                                **1 -> Online and available**                                 |
|                       |        |                             **2 -> Online but during the trip**                              |

---

### 1.1 **[GET]** api/v1/drivers/

It retrieves the drivers based on the request query parameters (if there is any). It also supports filtering the drivers based on **available_status** by putting the filtered condition in the request query parameters

#### Endpoint URL

```bash
http://localhost:80/server/api/v1/drivers/
```

#### Response

The returned results will be in the array of driver json object.

```json
[
  {
    "driver_id": "1",
    "first_name": "Run Lin",
    "last_name": "Xiong",
    "mobile_number": "6522222222",
    "email_address": "runlin@gmail.com",
    "identification_number": "T12345678A",
    "car_license_number": "h124j451k32jj123f",
    "available_status": 0
  },
  {
    "driver_id": "2",
    "first_name": "Zhi Quan",
    "last_name": "Henry Ong",
    "mobile_number": "6533333333",
    "email_address": "henryong@gmail.com",
    "identification_number": "T11111111C",
    "car_license_number": "dbaa541bcc85bcb3a8",
    "available_status": 1
  },
  {
    "driver_id": "3",
    "first_name": "Ming Han",
    "last_name": "Vincent Tee",
    "mobile_number": "6544444444",
    "email_address": "vincentminghan@gmail.com",
    "identification_number": "T22222222B",
    "car_license_number": "agfahudsi142kj42",
    "available_status": 2
  }
]
```

---

### 1.2 **[GET]** api/v1/drivers/:driverid

It retrieves the driver associated with the supplied driverID.

#### Endpoint URL

```bash
http://localhost:80/server/api/v1/drivers/1
```

#### Response

The returned results will be the driver json object. Return 404 if no record is found.

```json
{
  "driver_id": "1",
  "first_name": "Run Lin",
  "last_name": "Xiong",
  "mobile_number": "6522222222",
  "email_address": "runlin@gmail.com",
  "identification_number": "T12345678A",
  "car_license_number": "h124j451k32jj123f",
  "available_status": 0
}
```

---

### 1.3 **[POST]** api/v1/drivers/:driverid

It creates a driver in MySQL database by specific driverID. Information such as driver_id, first_name, last_name, mobile_number, email_address, identification_number, car_license_number must be supplied in the request body during registration.

#### Endpoint URL

```bash
http://localhost:80/server/api/v1/drivers/4
```

#### Request body

```json
{
  "driver_id": "4",
  "first_name": "hello",
  "last_name": "world",
  "mobile_number": "6544444444",
  "email_address": "helloworld@gmail.com",
  "identification_number": "T14124124A",
  "car_license_number": "asjkidh12hjdas"
}
```

#### Response

Case 1: If the compulsory driver information is not provided, it will return message which says the information is not correctly supplied<br>
Case 2: It will fail and return conflict status code if a driver with same driverID is already found in the database<br>
Case 3: Otherwise, it will return success message with status created code

```text
201 - Driver added: 4
409 - Duplicate driver ID
422 - Please supply driver information in JSON format
422 - The data in body and parameters do not match
```

---

### 1.4 **[PUT]** api/v1/drivers/:driverid

It is either used for creating or updating the driver depends whether the driverID exists. It allows updating fields for first_name, last_name, mobile_number, email_address, car_license_number, and available_status by putting the fields in the request body.

#### Endpoint URL

```bash
http://localhost:80/server/api/v1/drivers/4
```

#### Request body

```json
{
  "driver_id": "4",
  "mobile_number": "6522222222"
}
```

#### Response

Case 1: If driverID exists, update the driver using the information retrieved from request body<br>
Case 2: If driverID does not exist, create the driver using the information retrieved from request body

```text
201 - Driver added: 4
202 - Driver updated: 4
422 - Please supply driver information in JSON format
422 - The data in body and parameters do not match
```

## 2. Passenger

### Data Structure

|    Field Name    |  Type  |                                         Description                                          |
| :--------------: | :----: | :------------------------------------------------------------------------------------------: |
|   passenger_id   | string |                           The unique ID which identifies passenger                           |
|    first_name    | string |                               The first name of the passenger                                |
|    last_name     | string |                                The last name of the passenger                                |
|  mobile_number   | string |                              The mobile number of the passenger                              |
|  email_address   | string |                              The email address of the passenger                              |
| available_status |  int   |                        The available status number of the passenger.                         |
|                  |        | **0 -> used by golang to indicate whether the integer variable has been initialised or not** |
|                  |        |                                **1 -> Online and available**                                 |
|                  |        |                             **2 -> Online but during the trip**                              |

---

### 2.1 **[GET]** api/v1/passengers/

It retrieves the passengers from the database.

#### Endpoint URL

```bash
http://localhost:80/server/api/v1/passengers/
```

#### Response

The returned results will be in the array of passenger json object.

```json
[
  {
    "passenger_id": "1",
    "first_name": "Hou Man",
    "last_name": "Wai",
    "mobile_number": "655132123",
    "email_address": "hwendev@gmail.com",
    "available_status": 1
  },
  {
    "passenger_id": "2",
    "first_name": "Rui Quan",
    "last_name": "Zachary Hong",
    "mobile_number": "6512345678",
    "email_address": "zachary@gmail.com",
    "available_status": 1
  },
  {
    "passenger_id": "3",
    "first_name": "Yong Teng",
    "last_name": "Tee",
    "mobile_number": "6511111111",
    "email_address": "teeyongteng@gmail.com",
    "available_status": 1
  }
]
```

---

### 2.2 **[GET]** api/v1/passengers/:passengerid

It retrieves the passenger associated with the supplied passengerID.

#### Endpoint URL

```bash
http://localhost:80/server/api/v1/passengers/2
```

#### Response

The returned results will be the passenger json object. Return 404 if no record is found.

```json
{
  "passenger_id": "2",
  "first_name": "Rui Quan",
  "last_name": "Zachary Hong",
  "mobile_number": "6512345678",
  "email_address": "zachary@gmail.com",
  "available_status": 1
}
```

---

### 2.3 **[POST]** api/v1/passengers/:passengerid

It create a passenger in MySQL database by specific passengerID. Information such as passenger_id, first_name, last_name, mobile_number and email_address must be supplied in the request body during registration.

#### Endpoint URL

```bash
http://localhost:80/server/api/v1/passengers/4
```

#### Request body

```json
{
  "passenger_id": "4",
  "first_name": "hello",
  "last_name": "world",
  "mobile_number": "6512313215",
  "email_address": "helloworld@gmail.com"
}
```

#### Response

Case 1: If the compulsory passenger information is not provided, it will return message which says the information is not correctly supplied<br>
Case 2: It will fail and return conflict status code if a passenger with same passengerID is already found in the database<br>
Case 3: Otherwise, it will return success message with status created code

```text
201 - Passenger added: 4
409 - Duplicate passenger ID
422 - Please supply passenger information in JSON format
422 - The data in body and parameters do not match
```

---

### 2.4 **[PUT]** api/v1/passengers/:passengerid

It is either used for creating or updating the passenger depends whether the passengerID exists. It allows updating fields for first_name, last_name, mobile_number, email_address and available_status by putting the fields in the request body.

#### Endpoint URL

```bash
http://localhost:80/server/api/v1/passengers/4
```

#### Request body

```json
{
  "passenger_id": "4",
  "mobile_number": "655132123"
}
```

#### Response

Case 1: If passengerID exists, update the passenger using the information retrieved from request body<br>
Case 2: If passengerID does not exist, create the passenger using the information retrieved from request body

```text
201 - Passenger added: 4
202 - Passenger updated: 4
422 - Please supply passenger information in JSON format
422 - The data in body and parameters do not match
```

---

## 3. Trip Microservice

### Data Structure

|     Field Name      |   Type    |                                             Description                                             |
| :-----------------: | :-------: | :-------------------------------------------------------------------------------------------------: |
|       trip_id       |  string   |                                 The unique ID which identifies trip                                 |
|    passenger_id     |  string   |                              The unique ID which identifies passenger                               |
|      driver_id      |  string   |                                The unique ID which identifies driver                                |
| pickup_postal_code  |  string   |                                 The pickup postal code of the trip                                  |
| dropoff_postal_code |  string   |                                 The dropoff postal code of the trip                                 |
|    created_time     |  string   |                                    The created time of the trip                                     |
|   completed_time    |  string   |                       The completed time of the trip. Initially was set as 0.                       |
|      passenger      | Passenger | The passenger object of the trip. This field is used when returning results to the frontend server. |
|       driver        |  Driver   |  The driver object of the trip. This field is used when returning results to the frontend server.   |
|    trip_progress    |    int    |                                   The trip progress of the trip.                                    |
|                     |           |    **0 -> used by golang to indicate whether the integer variable has been initialised or not**     |
|                     |           |                **1 -> Created by passenger, but no driver is found to be assgined**                 |
|                     |           | **2 -> A driver was already assigned for the trip, but the driver has not inititated the trip yet** |
|                     |           |                     **3 -> The trip is ongoing (Driver has initiated the trip**                     |
|                     |           |                       **4 -> The trip has ended (Driver has ended the trip)**                       |

---

### 3.1 **[GET]** api/v1/trips/

It retrieves the trips based on the request parameters (if there is any) from the database. It also supports filtering the trips based on driver_id, passenger_id and trip_progress by putting the filtered condition in the request query parameters.

#### Endpoint URL

```bash
http://localhost:80/server/api/v1/trips?passenger_id=1
```

#### Response

The returned results will be in the array of trip json object.

```json
[
  {
    "trip_id": "1",
    "passenger_id": "1",
    "driver_id": "2",
    "pickup_postal_code": "642678",
    "dropoff_postal_code": "730022",
    "trip_progress": 1,
    "created_time": 1637424569024,
    "completed_time": 0,
    "passenger": {
      "passenger_id": "1",
      "first_name": "Hou Man",
      "last_name": "Wai",
      "mobile_number": "655132123",
      "email_address": "hwendev@gmail.com",
      "available_status": 1
    },
    "driver": {
      "driver_id": "2",
      "first_name": "Zhi Quan",
      "last_name": "Henry Ong",
      "mobile_number": "6533333333",
      "email_address": "henryong@gmail.com",
      "identification_number": "T11111111C",
      "car_license_number": "dbaa541bcc85bcb3a8",
      "available_status": 1
    }
  },
  {
    "trip_id": "3",
    "passenger_id": "1",
    "driver_id": "3",
    "pickup_postal_code": "333333",
    "dropoff_postal_code": "444444",
    "trip_progress": 1,
    "created_time": 1637424565554,
    "completed_time": 0,
    "passenger": {
      "passenger_id": "1",
      "first_name": "Hou Man",
      "last_name": "Wai",
      "mobile_number": "655132123",
      "email_address": "hwendev@gmail.com",
      "available_status": 1
    },
    "driver": {
      "driver_id": "3",
      "first_name": "Ming Han!!",
      "last_name": "Vincent Tee",
      "mobile_number": "6544444444",
      "email_address": "vincentminghan@gmail.com",
      "identification_number": "T22222222B",
      "car_license_number": "agfahudsi142kj42",
      "available_status": 2
    }
  }
]
```

---

### 3.2 **[GET]** api/v1/trips/:tripid

It retrieves the trip associated with the supplied tripID.

#### Endpoint URL

```bash
http://localhost:80/server/api/v1/trips/2
```

#### Response

The returned results will be the trip json object. Return 404 if no record is found.

```json
{
  "trip_id": "2",
  "passenger_id": "2",
  "driver_id": "3",
  "pickup_postal_code": "111111",
  "dropoff_postal_code": "222222",
  "trip_progress": 1,
  "created_time": 1637524569024,
  "completed_time": 0,
  "passenger": {
    "passenger_id": "2",
    "first_name": "Rui Quan",
    "last_name": "Zachary Hong",
    "mobile_number": "6512345678",
    "email_address": "zachary@gmail.com",
    "available_status": 1
  },
  "driver": {
    "driver_id": "3",
    "first_name": "Ming Han!!",
    "last_name": "Vincent Tee",
    "mobile_number": "6544444444",
    "email_address": "vincentminghan@gmail.com",
    "identification_number": "T22222222B",
    "car_license_number": "agfahudsi142kj42",
    "available_status": 2
  }
}
```

---

### 3.3 **[POST]** api/v1/trips/:tripid

It create a trip in MySQL database by specific tripID. Information such as trip_id, passenger_id, pickup_postal_code, dropoff_postal_code must be supplied in the request body during registration.

#### Endpoint URL

```bash
http://localhost:80/server/api/v1/passengers/4
```

#### Request body

```json
{
  "trip_id": "4",
  "passenger_id": "1",
  "pickup_postal_code": "642678",
  "dropoff_postal_code": "730022"
}
```

#### Response

Case 1: If the compulsory trip information is not provided, it will return message which says the information is not correctly supplied<br>
Case 2: It will fail and return conflict status code if a trip with same tripID is already found in the database<br>
Case 3: Otherwise, it will return success message with status created code

```text
201 - Trip added: 4
409 - Duplicate trip ID
422 - Please supply trip information in JSON format
422 - The data in body and parameters do not match
```

---

### 3.4 **[PUT]** api/v1/passengers/:passengerid

It is either used for creating or updating the trip depends whether the tripID exists. It allows updating fields for trip_progress and driver_id by putting the fields in the request body.

#### Endpoint URL

```bash
http://localhost:80/server/api/v1/trips/4
```

#### Request body

```json
{
  "trip_id": "4",
  "trip_progress": 2
}
```

#### Response

Case 1: If tripID exists, update the trip using the information retrieved from request body<br>
Case 2: If tripID does not exist, create the trip using the information retrieved from request body

```text
201 - Passenger added: 4
202 - Passenger updated: 4
422 - Please supply passenger information in JSON format
422 - The data in body and parameters do not match
```

---
