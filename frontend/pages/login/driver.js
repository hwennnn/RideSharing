import React, { useState } from 'react';
import { Input, Button } from 'semantic-ui-react'
import Head from 'next/head'
import Router from 'next/router';
import styles from '../../styles/Home.module.css'


export default function LoginDriver() {
    const [driverID, setDriverID] = useState(''); // '' is the initial state value

    function loginAsDriver() {
        if (driverID != '') {
            Router.push(`/driver/${driverID}`)
        }
    }

    return (
        <div className={styles.container}>
            <Head>
                <title>Login</title>
                <meta name="description" content="Generated by create next app" />
                <link rel="icon" href="/favicon.ico" />
            </Head>

            <main className={styles.main}>
                <h1 className={styles.title}>
                    Login as a driver
                </h1>

                <p className={styles.description}>
                    You will be prompted to enter your <span className={styles.blueColor}>driver id</span> in order to login.
                </p>

                <Input value={driverID} onChange={e => setDriverID(e.target.value)} focus placeholder='Enter your driver id' />
                <br />
                <Button onClick={loginAsDriver} primary>Login</Button>
            </main>

        </div>
    )
}

