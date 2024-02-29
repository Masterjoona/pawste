import { showToast } from "./toast.js";
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

    console.log(expiration);
    console.log(burn);
    console.log(password);
    console.log(syntax);
    console.log(privacy);
    console.log(file);

    const formData = new FormData();
    formData.append("expiration", expiration);
    formData.append("burn", burn);
    formData.append("password", password);
    formData.append("syntax", syntax);
    formData.append("privacy", privacy);
    formData.append("text", text);
    formData.append("file", file);

    submitFormData(formData);
    //showToast("info", "Your data is being submitted");
    //showToast("info", " and you will be redirected to the result page!");
}

let typedSequence = "";

// Add event listener to the document
document.addEventListener("keydown", function (event) {
    // Check if the event target is not the textarea with class "input-textarea"
    if (
        event.target.tagName !== "TEXTAREA" ||
        !event.target.classList.contains("input-textarea")
    ) {
        // Append the pressed key to the typed sequence
        typedSequence += event.key.toLowerCase();
        console.log(typedSequence);
        // Check if the typed sequence contains "neko"
        if (typedSequence.includes("neko")) {
            console.log('The sequence "neko" was typed');
            // Reset the typed sequence after logging
            typedSequence = "";
            showToast("info", "neko");
        }
    }
});

waitForElementToDisplay(
    "body > div.buttons > button.submit-button",
    function () {
        document
            .querySelector("body > div.buttons > button.submit-button")
            .addEventListener("click", submit);
    },
    500,
    5000,
);
