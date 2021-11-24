import { driverEndpointBaseURL } from '../config/baseURL';
import express from 'express';

const driverRouter = express.Router();

driverRouter.get("/", function (req, res) {
    console.log(`${driverEndpointBaseURL}`)
    res.redirect(`${driverEndpointBaseURL}`)
});

driverRouter.get("/{driverID}", function (req, res) {
    driverID = req.params['driverID']
    res.redirect(`${driverEndpointBaseURL}/${driverID}`)
});

export default driverRouter;