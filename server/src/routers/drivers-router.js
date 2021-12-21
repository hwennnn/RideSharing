import { driverEndpointBaseURL } from '../config/baseURL';
import express from 'express';
import axios from 'axios';
const url = require('url');

const driverRouter = express.Router();

// Redirect the requests to the driver microservice

driverRouter.get("/", async function (req, res) {
    const result = await axios.get(url.format({
        pathname: `${driverEndpointBaseURL}`,
        query: req.query,
    }));

    res.status(200).json(result.data);
});

driverRouter.get("/:driverID", async function (req, res) {
    let driverID = req.params.driverID;
    const result = await axios.get(`${driverEndpointBaseURL}/${driverID}`);

    res.status(200).json(result.data);
});

driverRouter.post("/:driverID", async function (req, res) {
    let driverID = req.params.driverID;
    const result = await axios.post(`${driverEndpointBaseURL}/${driverID}`);

    res.status(200).json(result.data);
});

driverRouter.put("/:driverID", async function (req, res) {
    let driverID = req.params.driverID;
    const result = await axios.put(`${driverEndpointBaseURL}/${driverID}`);

    res.status(200).json(result.data);
});

export default driverRouter;