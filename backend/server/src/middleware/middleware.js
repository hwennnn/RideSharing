export async function tokenAuthentication(req, res, next) {
    const authHeader = req.headers.authorization;

    const token = authHeader != undefined ? authHeader.split(' ')[1] : '';

    if (token == "2a6b36bf-61b9-4d0e-904c-7843e7b97308") {
        next();
    } else {
        res.sendStatus(403);
    }
}