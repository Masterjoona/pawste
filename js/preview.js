import { highlightElement } from "./speed-highlight/index.js";
import { detectLanguage } from "./speed-highlight/detect.js";
import { waitForElementToDisplay } from "./helpers.js";

export function languageOption(code) {
    const syntax = document.getElementById("syntax").value;
    if (syntax === "auto") {
        return detectLanguage(code);
    }
    if (syntax === "none") {
        return "";
    }
    return syntax;
}

function togglePreview() {
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
        previewTextArea.innerHTML = textarea.value;
        const language = languageOption(textarea.value);
        language === ""
            ? highlightElement(previewTextArea, undefined, "multiline")
            : highlightElement(previewTextArea, language, "multiline");
        previewTextArea.style.display = "inline-block";
        previewButton.textContent = "edit";
        themeBox.style.display = "block";
    }
}

waitForElementToDisplay(
    "body > div.buttons > button.preview-button",
    function () {
        document
            .querySelector("body > div.buttons > button.preview-button")
            .addEventListener("click", togglePreview);
    },
    500,
    5000,
);
