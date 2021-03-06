import axios from 'axios';
import { uuid } from 'uuidv4';
import { serverRequestBaseUrl, requestConfig, clientRequestBaseUrl } from './globals';

export async function getStaticPathForPassengers() {
    const response = await axios.get(`${serverRequestBaseUrl}/passengers/`, requestConfig);

    return response.data.map((passenger) => {
        return {
            params: {
                id: passenger.passenger_id
            }
        }
    })
}

export async function isPassengerExist(passengerID, sentFromClient = true) {
    try {
        const response = await axios.get(`${sentFromClient ? clientRequestBaseUrl : serverRequestBaseUrl}/passengers/${passengerID}`, requestConfig)
        console.log(response.status)
    } catch (err) {
        return false
    }

    return true
}

export async function getPassenger(passengerID) {
    const response = await axios.get(`${serverRequestBaseUrl}/passengers/${passengerID}`, requestConfig)

    return response.data
}

export async function retrieveCompletedTripsForPassenger(passengerID) {
    const response = await axios.get(`${serverRequestBaseUrl}/trips?passenger_id=${passengerID}&trip_progress=4`, requestConfig)

    return response.data
}

export async function createTripAsPassenger(passengerID, pickupPostalCode, dropoffPostalCode) {
    const trip_id = uuid()

    const body = {
        "trip_id": trip_id,
        "passenger_id": passengerID,
        "pickup_postal_code": pickupPostalCode,
        "dropoff_postal_code": dropoffPostalCode
    }

    const response = await axios.post(`${clientRequestBaseUrl}/trips/${trip_id}`, body, requestConfig)

    return response
}