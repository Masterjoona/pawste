function handleFiles(files) {
    console.log(files);
    // replace the button with the file name
    const filename = files[0].name;
    const fileButton = document.getElementById("attach-button");
    fileButton.innerHTML = "attached: " + filename;
    fileButton.classList.add("attached");
    // show the remove button
    const removeButton = document.getElementById("remove-button");
    removeButton.style.display = "block";
}

function removeFiles() {
    // remove the file name
    const fileButton = document.getElementById("attach-button");
    fileButton.innerHTML = "Attach";
    fileButton.classList.remove("attached");
    // hide the remove button
    const removeButton = document.getElementById("remove-button");
    removeButton.style.display = "none";
    // clear the input
    const fileInput = document.getElementById("file-input");
    fileInput.value = "";
}
