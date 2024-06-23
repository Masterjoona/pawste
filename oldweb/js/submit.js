import { waitForElementToDisplay } from "./helpers.js";
function submitFormData(formData) {
    fetch("/submit", {
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
            location.href = "/p/" + data.pasteUrl;
        })
        .catch((error) => {
            console.error("Error:", error);
            showToast("warning", "Failed to submit the data");
        });
}

function submit() {
    const expiration = document.getElementById("expiration").value;
    const burn = document.getElementById("burn").value;
    const password = document.getElementById("password").value;
    const syntax = document.getElementById("syntax").value;
    const privacy = document.getElementById("privacy").value;

    const text = document.getElementById("text-input").value;
    const file = document.getElementById("file-input").files[0] || null;

    console.log(expiration, burn, password, syntax, privacy, text, file);

    const formData = new FormData();
    formData.append("expiration", expiration);
    formData.append("burn", burn);
    formData.append("password", password);
    formData.append("syntax", syntax);
    formData.append("privacy", privacy);
    formData.append("text", text);
    formData.append("file", file);

    submitFormData(formData);
}

waitForElementToDisplay(
    "body > div.buttons > div.right-buttons > button.submit-button",
    function () {
        document
            .querySelector(
                "body > div.buttons > div.right-buttons > button.submit-button",
            )
            .addEventListener("click", submit);
    },
    500,
    5000,
);
