import { getStaticPathForPassengers } from '../../../utils/utils';

export async function getStaticProps({ params }) {
    // Add the "await" keyword like this:
    const id = params.id
    return {
        props: {
            id
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

export default function DriverHome({ id }) {
    return (
        'hello' + id
    )
}