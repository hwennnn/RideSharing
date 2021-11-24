import { driverEndpointBaseURL } from '../config/baseURL';
import express from 'express';

const driverRouter = express.Router();

driverRouter.get("/", function (req, res) {
    res.redirect(`${driverEndpointBaseURL}/`)
});

driverRouter.get("/:driverID", function (req, res) {
    let driverID = req.params.driverID
    res.redirect(307, `${driverEndpointBaseURL}/${driverID}`)
});

driverRouter.post("/:driverID", function (req, res) {
    let driverID = req.params.driverID
    res.redirect(307, `${driverEndpointBaseURL}/${driverID}`)
});

driverRouter.put("/:driverID", function (req, res) {
    let driverID = req.params.driverID
    res.redirect(307, `${driverEndpointBaseURL}/${driverID}`)
});

export default driverRouter;