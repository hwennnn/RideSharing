import { tripEndpointBaseURL } from '../config/baseURL';
import express from 'express';
const url = require('url');

const tripRouter = express.Router();

// Redirect the requests to the trip microservice
// HTTP 307 Temporary Redirect is used
// so that the method and the body of the original request are reused to perform the redirected request

tripRouter.get("/", function (req, res) {
    res.redirect(url.format({
        pathname: `${tripEndpointBaseURL}`,
        query: req.query,
    }))
});

tripRouter.get("/:tripID", function (req, res) {
    let tripID = req.params.tripID
    res.redirect(307, `${tripEndpointBaseURL}/${tripID}`)
});

tripRouter.post("/:tripID", function (req, res) {
    let tripID = req.params.tripID
    res.redirect(307, `${tripEndpointBaseURL}/${tripID}`)
});

tripRouter.put("/:tripID", function (req, res) {
    let tripID = req.params.tripID
    res.redirect(307, `${tripEndpointBaseURL}/${tripID}`)
});

export default tripRouter;