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