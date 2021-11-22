import React, { useState } from 'react'
import { Button, Input, Form } from 'semantic-ui-react'
import styles from '../../styles/Home.module.css'
import Head from 'next/head'

export default function RegisterDriver() {

    const [firstName, setFirstName] = useState('')
    const [lastName, setLastName] = useState('')
    const [mobileNumber, setMobileNumber] = useState('')
    const [emailAddress, setEmailAddress] = useState('')
    const [identificationNumber, setIdentificationNumber] = useState('')
    const [carLicenseNumber, setCarLicenseNumber] = useState('')

    function registerAsDriver() {
        if (firstName != '' && lastName != '' && mobileNumber != '' && emailAddress != '' && identificationNumber != '' && carLicenseNumber != '') {
            console.log(firstName, lastName, mobileNumber, emailAddress, identificationNumber, carLicenseNumber)
        }
    }


    return (
        <div className={styles.container}>
            <Head>
                <title>Login</title>
                <meta name="description" content="Generated by create next app" />
                <link rel="icon" href="/favicon.ico" />
            </Head>

            <h1 className={styles.title}>
                Register as a driver
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
                        <label>Identification Number</label>
                        <Input value={identificationNumber} onChange={e => setIdentificationNumber(e.target.value)} maxLength="30" fluid placeholder='Identification Number' />
                    </Form.Field>
                    <Form.Field>
                        <label>Car License Number</label>
                        <Input value={carLicenseNumber} onChange={e => setCarLicenseNumber(e.target.value)} maxLength="30" fluid placeholder='Car License Number' />
                    </Form.Field>
                </Form.Group>

                <Button onClick={registerAsDriver} type='submit'>Submit</Button>
            </Form>

        </div>
    );
}