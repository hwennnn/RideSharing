import { passengerEndpointBaseURL } from '../config/baseURL';
import express from 'express';

const passengerRouter = express.Router();

passengerRouter.get("/", function (req, res) {
    res.redirect(`${passengerEndpointBaseURL}/`)
});

passengerRouter.get("/:passengerID", function (req, res) {
    let passengerID = req.params.passengerID
    res.redirect(307, `${passengerEndpointBaseURL}/${passengerID}`)
});

passengerRouter.post("/:passengerID", function (req, res) {
    let passengerID = req.params.passengerID
    res.redirect(307, `${passengerEndpointBaseURL}/${passengerID}`)
});

passengerRouter.put("/:passengerID", function (req, res) {
    let passengerID = req.params.passengerID
    res.redirect(307, `${passengerEndpointBaseURL}/${passengerID}`)
});

export default passengerRouter;