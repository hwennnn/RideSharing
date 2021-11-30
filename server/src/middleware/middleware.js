export async function tokenAuthentication(req, res, next) {
    const authHeader = req.headers.authorization;

    const token = authHeader != undefined ? authHeader.split(' ')[1] : '';

    // the logic to authenticate token
    // Case 1: if the token is valid, allow the requests pass through the middleware
    // Case 2: Otherwise block it, and send 403 status code which indicates access forbidden
    if (token == "1467a2a8-fff7-45b5-986d-679382d0707a") {
        next();
    } else {
        res.sendStatus(403);
    }
}