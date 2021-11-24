import { tripEndpointBaseURL } from '../config/baseURL';
import express from 'express';
const url = require('url');

const tripRouter = express.Router();

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