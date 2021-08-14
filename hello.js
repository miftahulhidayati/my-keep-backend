function getDaysBetweenDates(date1, date2) {
    // Calculate the difference in milliseconds
    var difference = date2.getTime() - date1.getTime();
    // Convert the difference to days and return
    return Math.round(difference / (24 * 60 * 60 * 1000));
}