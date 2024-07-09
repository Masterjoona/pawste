import { toast } from "@zerodevx/svelte-toast";

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

export function viewFile(pastename, filename) {
    window.open("/p/" + pastename + "/f/" + filename);
}

export function timeDifference(timestamp) {
    const now = new Date();
    const target = new Date(timestamp);
    const diff = target - now;

    const seconds = Math.floor(diff / 1000);
    const minutes = Math.floor(seconds / 60);
    const hours = Math.floor(minutes / 60);
    const days = Math.floor(hours / 24);
    const weeks = Math.floor(days / 7);

    if (weeks > 0)
        return (
            weeks +
            (weeks === 1 ? " week" : " weeks") +
            (diff > 0 ? " from now" : " ago")
        );
    if (days > 0)
        return (
            days +
            (days === 1 ? " day" : " days") +
            (diff > 0 ? " from now" : " ago")
        );
    if (hours > 0)
        return (
            hours +
            (hours === 1 ? " hour" : " hours") +
            (diff > 0 ? " from now" : " ago")
        );
    if (minutes > 0)
        return (
            minutes +
            (minutes === 1 ? " minute" : " minutes") +
            (diff > 0 ? " from now" : " ago")
        );
    if (seconds > 0)
        return (
            seconds +
            (seconds === 1 ? " second" : " seconds") +
            (diff > 0 ? " from now" : " ago")
        );

    return "just now";
}

export function prettifyFileSize(size) {
    if (size < 1024) return size + " B";
    if (size < 1024 * 1024) return (size / 1024).toFixed(2) + " KB";
    if (size < 1024 * 1024 * 1024)
        return (size / (1024 * 1024)).toFixed(2) + " MB";
    return (size / (1024 * 1024 * 1024)).toFixed(2) + " GB";
}

export const successToast = (msg) => {
    toast.push(msg, {
        theme: {
            "--toastColor": "mintcream",
            "--toastBackground": "rgba(72,187,120,0.9)",
            "--toastBarBackground": "#2F855A",
        },
    });
};

export const failToast = (msg) => {
    toast.push(msg, {
        theme: {
            "--toastColor": "mintcream",
            "--toastBackground": "rgba(255,0,0,0.9)",
            "--toastBarBackground": "red",
        },
    });
};

export async function deletePaste(pasteName, successFunc) {
    const resp = await fetch(`/p/${pasteName}`, {
        method: "DELETE",
        body: JSON.stringify({ password }),
    });
    if (!resp.ok) {
        failToast("Failed to delete paste!");
    } else {
        successFunc();
    }
}
