import React, { useState } from 'react'
import { Button, Input, Form } from 'semantic-ui-react'
import styles from '../../styles/Home.module.css'
import Head from 'next/head'
import Router from 'next/router';
import { uuid } from 'uuidv4';
import axios from 'axios';
import { clientRequestBaseUrl, requestConfig } from '../../utils/globals';

export default function RegisterPassenger() {

    const [firstName, setFirstName] = useState('')
    const [lastName, setLastName] = useState('')
    const [mobileNumber, setMobileNumber] = useState('')
    const [emailAddress, setEmailAddress] = useState('')

    async function registerAsPassenger() {
        if (firstName != '' && lastName != '' && mobileNumber != '' && emailAddress != '') {
            console.log(firstName, lastName, mobileNumber, emailAddress)

            let passengerID = uuid()
            var body = {
                "passenger_id": passengerID,
                "first_name": firstName,
                "last_name": lastName,
                "mobile_number": mobileNumber,
                "email_address": emailAddress,
            }
            try {
                let response = await axios.post(`${clientRequestBaseUrl}/passengers/${passengerID}`, body, requestConfig);
                if (response.status == 201) {
                    console.log(response.data)
                    Router.push(`/passenger/${passengerID}`)
                } else {
                    alert(response.data)
                }
            } catch (err) {
                alert(err)
            }
        }
    }


    return (
        <div className={styles.container}>
            <Head>
                <title>Register</title>
                <meta name="description" content="Generated by create next app" />
                <link rel="icon" href="/favicon.ico" />
            </Head>

            <h1 className={styles.title}>
                Register as a passenger
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

                <Button onClick={registerAsPassenger} type='submit'>Submit</Button>
            </Form>

        </div>
    );
}