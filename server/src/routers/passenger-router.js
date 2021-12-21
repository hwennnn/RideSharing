import { passengerEndpointBaseURL } from '../config/baseURL';
import express from 'express';

const passengerRouter = express.Router();

// Redirect the requests to the passenger microservice

passengerRouter.get("/", async function (req, res) {
    const result = await axios.get(url.format({
        pathname: `${passengerEndpointBaseURL}`,
        query: req.query,
    }));

    res.status(200).json(result.data);
});

passengerRouter.get("/:passengerID", async function (req, res) {
    let passengerID = req.params.passengerID;
    const result = await axios.get(`${passengerEndpointBaseURL}/${passengerID}`);

    res.status(200).json(result.data);
});

passengerRouter.post("/:passengerID", async function (req, res) {
    let passengerID = req.params.passengerID;
    const result = await axios.post(`${passengerEndpointBaseURL}/${passengerID}`);

    res.status(200).json(result.data);
});

passengerRouter.put("/:passengerID", async function (req, res) {
    let passengerID = req.params.passengerID;
    const result = await axios.put(`${passengerEndpointBaseURL}/${passengerID}`);

    res.status(200).json(result.data);
});

export default passengerRouter;