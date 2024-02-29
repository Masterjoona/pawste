import {
    codeToHtml,
    bundledThemes,
    bundledLanguages,
} from "https://esm.sh/shiki@1.1.6";
import { waitForElementToDisplay } from "./helpers.js";

export function languageOption() {
    return document.getElementById("syntax").value;
}

async function togglePreview() {
    const previewButton = document.getElementById("preview-button");
    const previewTextArea = document.getElementById("preview");
    const textarea = document.getElementById("text-input");
    const themeBox = document.querySelector(
        "body > div.options > div:nth-child(6)",
    );

    if (textarea.style.display === "none") {
        textarea.style.display = "";
        previewTextArea.style.display = "none";
        previewButton.textContent = "preview";
        themeBox.style.display = "none";
    } else {
        textarea.style.display = "none";
        previewTextArea.innerHTML = await codeToHtml(textarea.value, {
            lang: languageOption(),
            theme: "solarized-dark",
        });
        previewTextArea.style.display = "inline-block";
        previewButton.textContent = "edit";
        themeBox.style.display = "block";
    }
}

function prettifyThemeName(theme) {
    return theme
        .split("-")
        .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
        .join(" ");
}

function addThemeOptions() {
    console.log(bundledThemes);
    const themeSelect = document.getElementById("theme");
    for (const theme in bundledThemes) {
        const option = document.createElement("option");
        option.value = theme;
        option.text = prettifyThemeName(theme);
        themeSelect.add(option);
    }
}

function addLanguageOptions() {
    console.log(bundledLanguages);
    const syntaxSelect = document.getElementById("syntax");
    console.log(Object.keys(bundledLanguages));
    for (const language in Object.keys(bundledLanguages)) {
        console.log("what the fuck");
        console.log(language);
    }
}

waitForElementToDisplay(
    "body > div.buttons > button.preview-button",
    function () {
        document
            .querySelector("body > div.buttons > button.preview-button")
            .addEventListener("click", async () => {
                await togglePreview();
            });
        addThemeOptions();
        //addLanguageOptions();
    },
    500,
    5000,
);
