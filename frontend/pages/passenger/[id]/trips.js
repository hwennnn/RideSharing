import { getStaticPathForPassengers, retrieveTripsForPassenger } from '../../../utils/passenger-utils';
import React from 'react'
import Head from 'next/head'
import { Table } from 'semantic-ui-react'
import styles from '../../../styles/Home.module.css'


export async function getStaticProps({ params }) {
    const passengerID = params.id
    const trips = await retrieveTripsForPassenger(passengerID);

    return {
        props: {
            trips
        }
    }
}

export async function getStaticPaths() {
    const paths = await getStaticPathForPassengers();

    return {
        paths,
        fallback: false
    }
}

export default function ViewTrips({ trips }) {

    const rows = trips != null ? trips.map(function (trip) {
        return (
            <Table.Row key={trip.trip_id}>
                <Table.Cell>{trip.trip_id}</Table.Cell>
                <Table.Cell>{`${trip.passenger.last_name} ${trip.passenger.first_name} #${trip.passenger.passenger_id}`}</Table.Cell>
                <Table.Cell>{`${trip.driver.last_name} ${trip.driver.first_name} #${trip.driver.driver_id}`}</Table.Cell>
                <Table.Cell>{trip.pickup_postal_code}</Table.Cell>
                <Table.Cell>{trip.dropoff_postal_code}</Table.Cell>
                <Table.Cell>{trip.created_time}</Table.Cell>
                <Table.Cell>{trip.completed_time}</Table.Cell>
            </Table.Row>
        )
    }) : 'There is no completed trip to view yet'

    return (

        <div className={styles.container}>
            <Head>
                <title>Completed Trips</title>
                <meta name="description" content="Generated by create next app" />
                <link rel="icon" href="/favicon.ico" />
            </Head>

            <h1 className={styles.title}>
                View Completed Trips
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
                    </Table.Row>
                </Table.Header>

                <Table.Body>
                    {rows}
                </Table.Body>

            </Table>

        </div>
    )
}