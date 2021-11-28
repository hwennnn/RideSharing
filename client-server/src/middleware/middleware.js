export async function tokenAuthentication(req, res, next) {
    const authHeader = req.headers.authorization;

    const token = authHeader != undefined ? authHeader.split(' ')[1] : '';

    if (token == "1467a2a8-fff7-45b5-986d-679382d0707a") {
        next();
    } else {
        res.sendStatus(403);
    }
}