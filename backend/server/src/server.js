import express from 'express';
import cors from "cors"
import helmet from "helmet"
import driversRouter from './routers/drivers-router';
import passengerRouter from './routers/passenger-router';
import tripRouter from './routers/trip-router';
import { tokenAuthentication } from './middleware/middleware';

const PORT = 4000;
const app = express();

app.use(helmet()); //safety
app.use(cors()); //safety
app.use(express.json()); //receive do respond with request

// the middleware to authenticate the token 
// to ensure the requests are valid and sent from the micorservices server.
app.use(tokenAuthentication)

app.use('/api/v1/drivers', driversRouter)
app.use('/api/v1/passengers', passengerRouter)
app.use('/api/v1/trips', tripRouter)


app.listen(PORT, async () => {
    console.log(`Listening on port ${PORT}`);
});