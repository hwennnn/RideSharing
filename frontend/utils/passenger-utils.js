import axios from 'axios';
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