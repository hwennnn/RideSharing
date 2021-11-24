import express from 'express';
import cors from "cors"
import helmet from "helmet"
import driversRouter from './routers/drivers-router';
import passengerRouter from './routers/passenger-router';

const PORT = 5000;
const app = express();

app.use(helmet()); //safety
app.use(cors()); //safety
app.use(express.json()); //receive do respond with request

app.get('/api/v1', function (req, res) {
    res.json('Hello World!')
})

app.use('/api/v1/drivers', driversRouter)
app.use('/api/v1/passengers', passengerRouter)


app.listen(PORT, async () => {
    console.log(`Listening on port ${PORT}`);
});