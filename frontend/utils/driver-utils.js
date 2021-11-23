import axios from 'axios';
import baseUrl from './baseUrl';

export async function getStaticPathForDrivers() {
    const response = await axios.get(`${baseUrl}/drivers/`);

    return response.data.map((driver) => {
        return {
            params: {
                id: driver.driver_id
            }
        }
    })
}

export async function getDriver(driverID) {
    const response = await axios.get(`${baseUrl}/drivers/${driverID}`)

    return response.data
}

export async function retrieveCompletedTripsForDriver(driverID) {
    const response = await axios.get(`${baseUrl}/trips?driver_id=${driverID}&trip_progress=3`)

    return response.data
}

export async function retrieveAvailableTripsForDriver() {
    const response = await axios.get(`${baseUrl}/trips?trip_progress=1`)

    return response.data
}

export async function initiateTripAsDriver(trip_id, driver_id) {
    const body = {
        "trip_id": trip_id,
        "driver_id": driver_id,
        "trip_progress": 2
    }
    console.log(body, `${baseUrl}/trips/${trip_id}`)
    const response = await axios.put(`${baseUrl}/trips/${trip_id}`, body)

    return response
}