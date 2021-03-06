import { getDriver, getStaticPathForDrivers } from '../../../utils/driver-utils';
import Head from 'next/head'
import Link from 'next/link'
import styles from '../../../styles/Home.module.css'

export async function getStaticProps({ params }) {
    const driverID = params.id
    const driver = await getDriver(driverID, false);

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

export default function DriverHome({ driver_id, last_name, first_name, available_status }) {
    let editProfileLink = `${driver_id}/edit`;
    let viewTripsLink = `${driver_id}/view_trips`
    let initiateTrip = `${driver_id}/initiate_trip`
    let endTrip = `${driver_id}/end_trip`

    return (
        <div className={styles.container}>
            <Head>
                <title>Home</title>
                <meta name="description" content="Generated by create next app" />
                <link rel="icon" href="/favicon.ico" />
            </Head>

            <main className={styles.main}>
                <h1 className={styles.title}>
                    U are now signed in as {last_name} {first_name} #{driver_id}
                </h1>

                {available_status == 2 &&
                    <p className={styles.description}>
                        You have been <span className={styles.blueColor}>assigned</span> a trip by the system
                    </p>
                }

                {available_status == 3 &&
                    <p className={styles.description}>
                        You are currently in an <span className={styles.blueColor}>ongoing</span> trip
                    </p>
                }

                <br />

                <div className={styles.grid}>
                    {available_status == 2 &&
                        <Link href={initiateTrip}>
                            <a className={styles.card}>
                                <h2>Initiate Trip &rarr;</h2>
                                <p>Choose and <span className={styles.blueColor}> initiate </span>a trip from the available list</p>
                            </a>
                        </Link>
                    }

                    {available_status == 3 &&
                        <Link href={endTrip}>
                            <a className={styles.card}>
                                <h2>End current ongoing trip &rarr;</h2>
                                <p>End the <span className={styles.blueColor}> current ongoing trip </span> if u have already reached the destination</p>
                            </a>
                        </Link>
                    }

                    <Link href={editProfileLink}>
                        <a className={styles.card}>
                            <h2>Edit driver profile &rarr;</h2>
                            <p>Edit your driver profile except <span className={styles.blueColor}> identification number </span></p>
                        </a>
                    </Link>

                    <Link href={viewTripsLink}>
                        <a className={styles.card}>
                            <h2>View past trips &rarr;</h2>
                            <p>View past completed trips in the table</p>
                        </a>
                    </Link>

                </div>
            </main>

        </div>
    )
}