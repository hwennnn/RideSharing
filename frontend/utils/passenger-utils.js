import axios from 'axios';
import { uuid } from 'uuidv4';
import baseUrl from './baseUrl';

export async function getStaticPathForPassengers() {
    const response = await axios.get(`${baseUrl}/passengers/`);

    return response.data.map((passenger) => {
        return {
            params: {
                id: passenger.passenger_id
            }
        }
    })
}

export async function getPassenger(passengerID) {
    const response = await axios.get(`${baseUrl}/passengers/${passengerID}`)

    return response.data
}

export async function retrieveTripsForPassenger(passengerID) {
    const response = await axios.get(`${baseUrl}/trips?passenger_id=${passengerID}`)

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

    const response = await axios.post(`${baseUrl}/trips/${trip_id}`, body)

    return response
}