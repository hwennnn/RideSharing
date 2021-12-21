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

    res.status(result.status).json(result.data);
});

tripRouter.get("/:tripID", async function (req, res) {
    let tripID = req.params.tripID;
    let body = req.body;
    const result = await axios.get(`${tripEndpointBaseURL}/${tripID}`, body);

    res.status(result.status).json(result.data);
});

tripRouter.post("/:tripID", async function (req, res) {
    let tripID = req.params.tripID;
    let body = req.body;
    const result = await axios.post(`${tripEndpointBaseURL}/${tripID}`, body);

    res.status(result.status).json(result.data);
});

tripRouter.put("/:tripID", async function (req, res) {
    let tripID = req.params.tripID;
    let body = req.body;
    const result = await axios.put(`${tripEndpointBaseURL}/${tripID}`, body);

    res.status(result.status).json(result.data);
});

export default tripRouter;