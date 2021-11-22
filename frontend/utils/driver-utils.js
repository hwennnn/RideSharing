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