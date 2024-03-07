import { showToast } from "./toast.js";
import { waitForElementToDisplay } from "./helpers.js";
import { languageOption } from "./preview.js";
import { codeToHtml } from "https://esm.sh/shiki@1.1.7";

async function changeTheme() {
    const theme = document.getElementById("theme").value;
    document.getElementById("preview").innerHTML = await codeToHtml(
        document.getElementById("text-input").value,
        {
            lang: languageOption(),
            theme: theme,
        },
    );
    showToast("info", `Theme changed to ${theme}`);
}

waitForElementToDisplay(
    "#theme",
    function () {
        document
            .getElementById("theme")
            .addEventListener("change", async () => {
                changeTheme();
            });
    },
    500,
    5000,
);
