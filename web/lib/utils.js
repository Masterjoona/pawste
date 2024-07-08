export function truncateFilename(filename, maxLength = 30) {
    const extIndex = filename.lastIndexOf(".");
    const name = filename.substring(0, extIndex);
    const ext = filename.substring(extIndex);

    if (name.length + ext.length <= maxLength) {
        return filename;
    }

    const charsToShow = maxLength - ext.length - 3;
    const startChars = Math.ceil(charsToShow / 2);
    const endChars = Math.floor(charsToShow / 2);

    return (
        name.substring(0, startChars) +
        "..." +
        name.substring(name.length - endChars) +
        ext
    );
}

export function viewFile(filename) {
    window.open("/p/" + paste.PasteName + "/f/" + filename);
}
