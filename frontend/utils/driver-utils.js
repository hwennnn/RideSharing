import axios from 'axios';
import { baseUrl, requestConfig } from './globals';

export async function getStaticPathForDrivers() {
    const response = await axios.get(`${baseUrl}/drivers/`, requestConfig);

    return response.data.map((driver) => {
        return {
            params: {
                id: driver.driver_id
            }
        }
    })
}

export async function getDriver(driverID) {
    const response = await axios.get(`${baseUrl}/drivers/${driverID}`, requestConfig)

    return response.data
}

export async function retrieveCompletedTripsForDriver(driverID) {
    const response = await axios.get(`${baseUrl}/trips?driver_id=${driverID}&trip_progress=3`, requestConfig)

    return response.data
}

export async function retrieveAvailableTripsForDriver() {
    const response = await axios.get(`${baseUrl}/trips?trip_progress=1`, requestConfig)

    return response.data
}

export async function initiateTripAsDriver(trip_id, driver_id) {
    const body = {
        "trip_id": trip_id,
        "driver_id": driver_id,
        "trip_progress": 2
    }

    const response = await axios.put(`${baseUrl}/trips/${trip_id}`, body, requestConfig)

    return response
}

export async function retrieveOngoingTripForDriver(driverID) {
    const response = await axios.get(`${baseUrl}/trips?driver_id=${driverID}&trip_progress=2`, requestConfig)

    return (response.data)[0]
}

export async function endTripAsDriver(trip_id) {
    const body = {
        "trip_id": trip_id,
        "trip_progress": 3
    }

    const response = await axios.put(`${baseUrl}/trips/${trip_id}`, body, requestConfig)

    return response
}