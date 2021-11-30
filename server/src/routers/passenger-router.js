import { passengerEndpointBaseURL } from '../config/baseURL';
import express from 'express';

const passengerRouter = express.Router();

// Redirect the requests to the passenger microservice
// HTTP 307 Temporary Redirect is used
// so that the method and the body of the original request are reused to perform the redirected request

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