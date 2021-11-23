import { getDriver, getStaticPathForDrivers } from '../../../utils/driver-utils';
import React, { useState } from 'react'
import { Button, Input, Form } from 'semantic-ui-react'
import styles from '../../../styles/Home.module.css'
import Head from 'next/head'
import Router from 'next/router';
import axios from 'axios';
import { baseUrl, requestConfig } from '../../../utils/globals';


export async function getStaticProps({ params }) {
    const driverID = params.id
    const driver = await getDriver(driverID);
    return {
        props: {
            ...driver
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

export default function EditDriver({ driver_id, first_name, last_name, mobile_number, email_address, car_license_number }) {
    const [firstName, setFirstName] = useState(first_name)
    const [lastName, setLastName] = useState(last_name)
    const [mobileNumber, setMobileNumber] = useState(mobile_number)
    const [emailAddress, setEmailAddress] = useState(email_address)
    const [carLicenseNumber, setCarLicenseNumber] = useState(car_license_number)

    async function updateAsDriver() {
        if (firstName != '' && lastName != '' && mobileNumber != '' && emailAddress != '' && carLicenseNumber != '') {
            var body = {
                "driver_id": driver_id,
                "first_name": firstName,
                "last_name": lastName,
                "mobile_number": mobileNumber,
                "email_address": emailAddress,
                "car_license_number": carLicenseNumber
            }
            try {
                let response = await axios.put(`${baseUrl}/drivers/${driver_id}`, body, requestConfig);

                if (response.status == 202) {
                    console.log(response.data)
                    Router.push(`/driver/${driver_id}`)
                } else {
                    alert(response.data)
                }
            } catch (err) {
                alert(err)
            }
        }
    }

    function backToDriverHome() {
        Router.push(`/driver/${driver_id}`)
    }

    return (
        <div className={styles.container}>
            <Head>
                <title>Edit Driver Profile</title>
                <meta name="description" content="Generated by create next app" />
                <link rel="icon" href="/favicon.ico" />
            </Head>

            <h1 className={styles.title}>
                Edit Driver Profile
            </h1>

            <br />

            <Form>
                <Form.Group widths='equal'>
                    <Form.Field>
                        <label>First Name</label>
                        <Input value={firstName} onChange={e => setFirstName(e.target.value)} maxLength="30" fluid placeholder='First Name' />
                    </Form.Field>
                    <Form.Field>
                        <label>Last name</label>
                        <Input value={lastName} onChange={e => setLastName(e.target.value)} maxLength="30" fluid placeholder='Last Name' />
                    </Form.Field>
                </Form.Group>

                <Form.Group widths='equal'>
                    <Form.Field>
                        <label>Mobile Number</label>
                        <Input value={mobileNumber} onChange={e => setMobileNumber(e.target.value)} maxLength="15" fluid placeholder='Mobile Number' />
                    </Form.Field>
                    <Form.Field>
                        <label>Email Address</label>
                        <Input value={emailAddress} onChange={e => setEmailAddress(e.target.value)} maxLength="30" type="email" fluid placeholder='Email Address' />
                    </Form.Field>
                </Form.Group>

                <Form.Group widths='equal'>
                    <Form.Field>
                        <label>Car License Number</label>
                        <Input value={carLicenseNumber} onChange={e => setCarLicenseNumber(e.target.value)} maxLength="30" fluid placeholder='Car License Number' />
                    </Form.Field>
                </Form.Group>

                <Button onClick={updateAsDriver} type='submit'>Submit</Button>
            </Form>

            <br />
            <Button primary onClick={backToDriverHome} type='submit'>Back To Home</Button>
        </div>
    )
}