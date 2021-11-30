import { driverEndpointBaseURL } from '../config/baseURL';
import express from 'express';
const url = require('url');

const driverRouter = express.Router();

// Redirect the requests to the driver microservice
// HTTP 307 Temporary Redirect is used
// so that the method and the body of the original request are reused to perform the redirected request

driverRouter.get("/", function (req, res) {
    res.redirect(url.format({
        pathname: `${driverEndpointBaseURL}`,
        query: req.query,
    }))
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