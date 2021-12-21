import axios from 'axios';
import { serverRequestBaseUrl, requestConfig, clientRequestBaseUrl } from './globals';

export async function getStaticPathForDrivers() {
    const response = await axios.get(`${serverRequestBaseUrl}/drivers/`, requestConfig);

    return response.data.map((driver) => {
        return {
            params: {
                id: driver.driver_id
            }
        }
    })
}

export async function isDriverExist(driverID, sentFromClient = true) {
    try {
        const response = await axios.get(`${sentFromClient ? clientRequestBaseUrl : serverRequestBaseUrl}/drivers/${driverID}`, requestConfig)
        console.log(response.status)
    } catch (err) {
        return false
    }

    return true
}

export async function getDriver(driverID) {
    const response = await axios.get(`${serverRequestBaseUrl}/drivers/${driverID}`, requestConfig)

    return response.data
}

export async function retrieveCompletedTripsForDriver(driverID) {
    const response = await axios.get(`${serverRequestBaseUrl}/trips?driver_id=${driverID}&trip_progress=4`, requestConfig)

    return response.data
}

export async function retrieveAvailableTripsForDriver() {
    const response = await axios.get(`${serverRequestBaseUrl}/trips?trip_progress=2`, requestConfig)

    return response.data
}

export async function initiateTripAsDriver(trip_id, driver_id) {
    const body = {
        "trip_id": trip_id,
        "driver_id": driver_id,
        "trip_progress": 3
    }

    const response = await axios.put(`${clientRequestBaseUrl}/trips/${trip_id}`, body, requestConfig)

    return response
}

export async function retrieveOngoingTripForDriver(driverID) {
    const response = await axios.get(`${serverRequestBaseUrl}/trips?driver_id=${driverID}&trip_progress=3`, requestConfig)
    console.log(response);
    return (response.data)[0]
}

export async function endTripAsDriver(trip_id) {
    const body = {
        "trip_id": trip_id,
        "trip_progress": 4
    }

    const response = await axios.put(`${clientRequestBaseUrl}/trips/${trip_id}`, body, requestConfig)

    return response
}