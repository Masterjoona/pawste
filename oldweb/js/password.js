// for node
// const crypto = require("crypto");
import { waitForElementToDisplay } from "./helpers.js";

async function HashPassword(password) {
    // Convert the password string to an array buffer
    const passwordBuffer = new TextEncoder().encode(password);

    try {
        // Generate a cryptographic hash using the SHA-256 algorithm
        const hashBuffer = await crypto.subtle.digest(
            "SHA-256",
            passwordBuffer,
        );

        // Convert the hash buffer to a hexadecimal string
        const hashArray = Array.from(new Uint8Array(hashBuffer));
        const hashHex = hashArray
            .map((byte) => byte.toString(16).padStart(2, "0"))
            .join("");

        return hashHex;
    } catch (error) {
        console.error("Error hashing password:", error);
        return null;
    }
}

/*
// Example usage:
const password = "test";
HashPassword(password).then((hash) => {
    console.log("Hashed password:", hash);
});
*/

function togglePasswordBoxVisibility() {
    const privacy = document.querySelector(
        "body > div.options > div:nth-child(5)",
    );
    const privacySelect = document.getElementById("privacy");
    if (privacySelect.value === "private" || privacySelect.value === "secret") {
        privacy.style.display = "block";
    } else {
        privacy.style.display = "none";
    }
}

waitForElementToDisplay(
    "body > div.options > div:nth-child(4)",
    function () {
        document
            .getElementById("privacy")
            .addEventListener("change", togglePasswordBoxVisibility);
    },
    500,
    5000,
);
