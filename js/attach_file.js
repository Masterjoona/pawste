window.addEventListener("paste", (e) => {
    if (e.clipboardData.files.length === 0) return;
    const fileInput = document.getElementById("file-input");
    console.log(e.clipboardData.files);
    fileInput.files = e.clipboardData.files;
    const event = new Event("change");
    fileInput.dispatchEvent(event);
});

function handleFiles(files) {
    const filename = files[0].name;
    const fileButton = document.getElementById("attach-button");
    let prettyFilename = filename;
    if (filename.length > 10) {
        prettyFilename = filename.substr(0, 5) + "..." + filename.substr(-5);
    }
    fileButton.innerHTML = "attached: " + prettyFilename;
    fileButton.classList.add("attached");

    const removeButton = document.getElementById("remove-button");
    removeButton.style.display = "block";
}

function removeFiles() {
    const fileButton = document.getElementById("attach-button");
    fileButton.innerHTML = "Attach";
    fileButton.classList.remove("attached");

    const removeButton = document.getElementById("remove-button");
    removeButton.style.display = "none";

    const fileInput = document.getElementById("file-input");
    fileInput.value = "";
    const event = new Event("change");
    fileInput.dispatchEvent(event);
}
