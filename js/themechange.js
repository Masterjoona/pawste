import { showToast } from "./toast.js";
import { waitForElementToDisplay } from "./helpers.js";

function changeTheme() {
    const theme = document.getElementById("theme").value;
    const themeStyle = document.getElementById("theme-stylesheet");
    const previewButton = document.getElementById("preview-button");
    if (previewButton.textContent === "preview") {
        showToast(
            "warning",
            "You need to preview the text to see the theme change"
        );
        return;
    }
    themeStyle.href = `css/themes/${theme}.css`;
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
    5000
);

document.getElementById("theme").addEventListener("change", changeTheme);
