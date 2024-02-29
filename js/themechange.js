import { showToast } from "./toast.js";
import { waitForElementToDisplay } from "./helpers.js";

function changeTheme() {
    const theme = document.getElementById("theme").value;

    showToast("info", `Theme changed to ${theme}`);
}

waitForElementToDisplay(
    "#theme",
    function () {
        document
            .getElementById("theme")
            .addEventListener("change", changeTheme);
    },
    500,
    5000,
);

document.getElementById("theme").addEventListener("change", changeTheme);
