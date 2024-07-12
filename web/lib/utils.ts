import { toast } from "@zerodevx/svelte-toast";
import { FileDb, FileType } from "./types";

export function truncateFilename(filename: string, maxLength = 30) {
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

export function viewFile(pastename: string, filename: string) {
    window.open("/p/" + pastename + "/f/" + filename);
}

export function isFileDb(file: FileType): file is FileDb {
    return (file as FileDb)?.Name !== undefined;
}

export function timeDifference(timestamp: number) {
    timestamp *= 1000;

    const now = new Date();
    const target = new Date(timestamp);
    const diff = target.getTime() - now.getTime();

    const absDiff = Math.abs(diff);

    const seconds = Math.floor(absDiff / 1000);
    const minutes = Math.floor(seconds / 60);
    const hours = Math.floor(minutes / 60);
    const days = Math.floor(hours / 24);
    const weeks = Math.floor(days / 7);

    let timeUnit: string;
    let timeValue: number;

    if (weeks > 0) {
        timeUnit = weeks === 1 ? "week" : "weeks";
        timeValue = weeks;
    } else if (days > 0) {
        timeUnit = days === 1 ? "day" : "days";
        timeValue = days;
    } else if (hours > 0) {
        timeUnit = hours === 1 ? "hour" : "hours";
        timeValue = hours;
    } else if (minutes > 0) {
        timeUnit = minutes === 1 ? "minute" : "minutes";
        timeValue = minutes;
    } else if (seconds > 0) {
        timeUnit = seconds === 1 ? "second" : "seconds";
        timeValue = seconds;
    } else {
        return "just now";
    }

    const suffix = diff > 0 ? "" : "ago";

    return `${timeValue} ${timeUnit} ${suffix}`;
}

export function prettifyFileSize(size: number) {
    if (size < 1024) return size + " B";
    if (size < 1024 * 1024) return (size / 1024).toFixed(2) + " KB";
    if (size < 1024 * 1024 * 1024)
        return (size / (1024 * 1024)).toFixed(2) + " MB";
    return (size / (1024 * 1024 * 1024)).toFixed(2) + " GB";
}

export async function deletePaste(
    pasteName: string,
    password: string,
    onSuccess: () => void,
) {
    const resp = await fetch(`/p/${pasteName}`, {
        method: "DELETE",
        headers: {
            "Content-Type": "application/json",
            password,
        },
    });
    if (!resp.ok) {
        failToast("Failed to delete paste!");
        return;
    }
    onSuccess();
}

export const successToast = (msg: string) => {
    toast.push(msg, {
        theme: {
            "--toastColor": "mintcream",
            "--toastBackground": "rgba(72,187,120,0.9)",
            "--toastBarBackground": "#2F855A",
        },
    });
};

export const failToast = (msg: string) => {
    toast.push(msg, {
        theme: {
            "--toastColor": "mintcream",
            "--toastBackground": "rgba(255,0,0,0.9)",
            "--toastBarBackground": "red",
        },
    });
};
