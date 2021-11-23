import dateFormat from "dateformat";

export function formatDateStringFromMs(ms) {
    var dt = new Date(ms);
    return dateFormat(dt, "dS mmmm yyyy, h:MM:ss TT");
}