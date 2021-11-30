export async function tokenAuthentication(req, res, next) {
    const authHeader = req.headers.authorization;

    const token = authHeader != undefined ? authHeader.split(' ')[1] : '';

    // the logic to authenticate token
    // Case 1: if the token is valid, allow the requests pass through the middleware
    // Case 2: Otherwise block it, and send 403 status code which indicates access forbidden
    if (token == "2a6b36bf-61b9-4d0e-904c-7843e7b97308") {
        next();
    } else {
        res.sendStatus(403);
    }
}