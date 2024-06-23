import { showToast } from "./toast.js";
import { waitForElementToDisplay } from "./helpers.js";

function reloadConfig() {
    const formData = new FormData();
    formData.append("password", document.getElementById("password").value);
    fetch("/admin/reload-config", {
        method: "POST",
        body: formData,
    })
        .then((response) => {
            if (!response.ok) {
                throw new Error(
                    `Error ${response.status}: ${response.statusText}`,
                );
            }
            return response.json();
        })
        .then((data) => {
            console.log(data);
            showToast("success", "Config reloaded");
        })
        .catch((error) => {
            console.error("Error:", error);
            showToast("warning", "Failed to reload the config");
        });
}

waitForElementToDisplay(
    "#reloadconfig",
    () => {
        document
            .getElementById("reloadconfig")
            .addEventListener("click", reloadConfig);
    },
    500,
    5000,
);
