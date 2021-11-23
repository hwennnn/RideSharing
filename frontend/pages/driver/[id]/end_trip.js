import { endTripAsDriver, getStaticPathForDrivers, retrieveOngoingTripForDriver } from '../../../utils/driver-utils';
import React from 'react'
import Head from 'next/head'
import { Table, Button } from 'semantic-ui-react'
import styles from '../../../styles/Home.module.css'
import Router from 'next/router';
import { formatDateStringFromMs } from '../../../utils/date-utils'


export async function getStaticProps({ params }) {
    const driverID = params.id
    const trip = await retrieveOngoingTripForDriver(driverID);

    return {
        props: {
            trip
        }
    }
}

export async function getStaticPaths() {
    const paths = await getStaticPathForDrivers();

    return {
        paths,
        fallback: false
    }
}

export default function EndTrip({ trip }) {

    async function endTrip() {
        let response = await endTripAsDriver(trip.trip_id)
        if (response.status == 202) {
            Router.push(`/driver/${trip.driver.driver_id}`)
        } else {
            alert(response.data)
        }
    }

    return (

        <div className={styles.container}>
            <Head>
                <title>End Trip</title>
                <meta name="description" content="Generated by create next app" />
                <link rel="icon" href="/favicon.ico" />
            </Head>

            <h1 className={styles.title}>
                End current ongoing trip
            </h1>

            <br />

            <Table celled>
                <Table.Header>
                    <Table.Row>
                        <Table.HeaderCell>TripID</Table.HeaderCell>
                        <Table.HeaderCell>Passenger</Table.HeaderCell>
                        <Table.HeaderCell>Driver</Table.HeaderCell>
                        <Table.HeaderCell>PickupPostalCode</Table.HeaderCell>
                        <Table.HeaderCell>DropoffPostalCode</Table.HeaderCell>
                        <Table.HeaderCell>CreatedTime</Table.HeaderCell>
                        <Table.HeaderCell>CompletedTime</Table.HeaderCell>
                        <Table.HeaderCell>Actions</Table.HeaderCell>
                    </Table.Row>
                </Table.Header>

                <Table.Body>
                    <Table.Row key={trip.trip_id}>
                        <Table.Cell>{trip.trip_id}</Table.Cell>
                        <Table.Cell>{`${trip.passenger.last_name} ${trip.passenger.first_name} #${trip.passenger.passenger_id}`}</Table.Cell>
                        <Table.Cell>{`${trip.driver.last_name} ${trip.driver.first_name} #${trip.driver.driver_id}`}</Table.Cell>
                        <Table.Cell>{trip.pickup_postal_code}</Table.Cell>
                        <Table.Cell>{trip.dropoff_postal_code}</Table.Cell>
                        <Table.Cell>{formatDateStringFromMs(trip.created_time)}</Table.Cell>
                        <Table.Cell>{formatDateStringFromMs(trip.completed_time)}</Table.Cell>
                        <Table.Cell><Button onClick={() => endTrip()} >End Trip</Button></Table.Cell>
                    </Table.Row>
                </Table.Body>

            </Table>

        </div>
    )
}