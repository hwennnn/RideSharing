import { tripEndpointBaseURL } from '../config/baseURL';
import express from 'express';
import axios from 'axios';
const url = require('url');

const tripRouter = express.Router();

// Redirect the requests to the trip microservice

tripRouter.get("/", async function (req, res) {
    const result = await axios.get(url.format({
        pathname: `${tripEndpointBaseURL}`,
        query: req.query,
    }));

    res.status(200).json(result.data);
});

tripRouter.get("/:tripID", async function (req, res) {
    let tripID = req.params.tripID;
    const result = await axios.get(`${tripEndpointBaseURL}/${tripID}`);

    res.status(200).json(result.data);
});

tripRouter.post("/:tripID", async function (req, res) {
    let tripID = req.params.tripID;
    const result = await axios.post(`${tripEndpointBaseURL}/${tripID}`);

    res.status(200).json(result.data);
});

tripRouter.put("/:tripID", async function (req, res) {
    let tripID = req.params.tripID;
    const result = await axios.put(`${tripEndpointBaseURL}/${tripID}`);

    res.status(200).json(result.data);
});

export default tripRouter;