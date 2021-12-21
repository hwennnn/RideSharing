import { passengerEndpointBaseURL } from '../config/baseURL';
import express from 'express';
import axios from 'axios';

const passengerRouter = express.Router();

// Redirect the requests to the passenger microservice

passengerRouter.get("/", async function (req, res) {
    const result = await axios.get(`${passengerEndpointBaseURL}/`);

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