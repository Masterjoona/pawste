export function showToast(type, message) {
    const toastContainer = document.getElementById("toastContainer");
    const existingToasts = toastContainer.querySelectorAll(".toast");
    if (existingToasts.length > 0) {
        moveToastUp(existingToasts);
    }
    const toast = document.createElement("div");
    toast.className = "toast";

    const span = document.createElement("span");
    span.textContent = message;
    toast.appendChild(span);

    toastContainer.appendChild(toast);

    if (type === "warning") {
        toast.classList.add("warning-toast");
    } else if (type === "info") {
        toast.classList.add("info-toast");
    } else {
        toast.classList.add("default-toast");
    }

    toastContainer.appendChild(toast);
    toast.offsetHeight;
    toast.classList.add("show-toast");

    setTimeout(() => {
        toast.classList.remove("show-toast");
        setTimeout(() => {
            toast.remove();
            if (!toastContainer.querySelector(".toast")) {
                resetToastPositions(existingToasts);
            }
        }, 500);
    }, 3000);
}

function moveToastUp(toasts) {
    toasts.forEach((toast) => {
        const newBottom =
            parseInt(window.getComputedStyle(toast).bottom) +
            toast.offsetHeight +
            10;
        toast.style.bottom = newBottom + "px";
    });
}

function resetToastPositions(toasts) {
    toasts.forEach((toast) => {
        toast.style.bottom = "20px";
    });
}
