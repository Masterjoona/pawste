<script>
    import { truncateFilename, viewFile } from "../lib/utils.js";
    import { toast } from "@zerodevx/svelte-toast";

    export let paste;
    export let files;

    let removedFiles = [];
    let newFiles = [];
    let imageSources = [];
    let newContent = paste.Content;

    function handleAttachFiles(event) {
        const files = event.target.files;
        for (let file of files) {
            newFiles = [...newFiles, file];
            if (file.type.startsWith("image/")) {
                const reader = new FileReader();
                reader.onload = (e) => {
                    imageSources = [...imageSources, e.target.result];
                };
                reader.readAsDataURL(file);
            } else {
                imageSources = [...imageSources, null];
            }
        }
    }

    function removeOldFile(filename) {
        removedFiles = [...removedFiles, filename];
        files = files.filter((file) => file.Name !== filename);
    }
    function removeNewFile(filename) {
        newFiles = newFiles.filter((file) => file.name !== filename);
    }
    function handleSave() {
        const noNewFiles = !newFiles.length;
        console.log(files, removedFiles, newFiles, newContent);
        if (!files.length && noNewFiles && !newContent) {
            toast.push("You must provide content or attach files!", {
                theme: {
                    "--toastColor": "mintcream",
                    "--toastBackground": "rgba(255,0,0,0.9)",
                    "--toastBarBackground": "red",
                },
            });
            return;
        }
        if (
            newContent === paste.Content &&
            noNewFiles &&
            !removedFiles.length
        ) {
            toast.push("No changes detected!", {
                theme: {
                    "--toastColor": "mintcream",
                    "--toastBackground": "rgba(255,0,0,0.9)",
                    "--toastBarBackground": "red",
                },
            });
            return;
        }
        toast.push("Changes saved!", {
            theme: {
                "--toastColor": "mintcream",
                "--toastBackground": "rgba(72,187,120,0.9)",
                "--toastBarBackground": "#2F855A",
            },
        });
    }
</script>

<div id="container">
    <div class="card">
        <div class="properties">
            <h>{paste.PasteName}</h>
            <div class="spacer"></div>
            <div class="icon-container">
                <p>{paste.ReadCount} <i class="fa-solid fa-eye"></i></p>
                <p>
                    {paste.Content.length}
                    <i class="fa-solid fa-file-lines"></i>
                </p>
                <p>{paste.Expire} <i class="fa-solid fa-clock"></i></p>
            </div>
        </div>
        <textarea bind:value={newContent}></textarea>
        <div class="buttons">
            <input
                type="file"
                multiple
                on:change={handleAttachFiles}
                style="display: none;"
                id="file-input" />
            <button
                on:click={() => document.getElementById("file-input").click()}
                >Attach Files</button>
            <button on:click={handleSave}>Save</button>
        </div>
        <div class="file-list">
            <p>Current Files:</p>
            {#each files as file}
                <div class="file-item">
                    <span
                        >{truncateFilename(file.Name)} - {(
                            file.Size / 1024
                        ).toFixed(2)} KB</span>
                    <button on:click={() => removeOldFile(file.Name)}
                        >Remove</button>
                    <button on:click={() => viewFile(file.Name)}>View</button>
                </div>
            {/each}
            <p>New Files:</p>
            {#each newFiles as file, index}
                <div class="file-item">
                    {#if file.type.startsWith("image/")}
                        <img
                            src={imageSources[index]}
                            alt={file.name}
                            class="thumbnail" />
                    {/if}
                    <span
                        >{truncateFilename(file.name)} - {(
                            file.size / 1024
                        ).toFixed(2)} KB</span>
                    <button on:click={() => removeNewFile(file.name)}
                        >Remove</button>
                </div>
            {/each}
        </div>
    </div>
</div>

<style>
    :root {
        --font-size: 1.2em;
    }

    #container {
        height: 100%;
        width: 100%;
        display: flex;
        justify-content: center;
        align-items: center;
        font-family: var(--main-font);
    }

    .card {
        width: 90%;
        display: flex;
        justify-content: center;
        align-items: center;
        flex-direction: column;
        background-color: #2a2a2a;
        border-radius: 10px;
        padding: 16px;
    }

    .properties {
        width: 100%;
        display: flex;
        flex-direction: row;
        justify-content: space-evenly;
        margin-top: 0px;
    }
    .properties {
        width: 100%;
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        margin-top: 0px;
    }

    .spacer {
        flex-grow: 1;
    }

    .icon-container {
        display: flex;
        flex-direction: row;
        gap: 10px;
    }

    textarea {
        width: 99%;
        height: 200px;
        font-family: var(--code-font);
        background-color: #1b1b22;
        color: white;
        border: none;
        border-radius: 10px;
        padding: 10px;
        resize: vertical;
        margin-bottom: 10px;
    }

    .buttons {
        display: flex;
        width: 100%;
        justify-content: space-evenly;
    }

    button {
        background-color: var(--main-color);
        color: white;
        border: none;
        padding: 10px 10px;
        border-radius: 5px;
        cursor: pointer;
        font-family: var(--main-font);
        font-size: var(--font-size);
    }

    button:hover {
        background-color: var(--main-color-dark);
    }

    .file-list {
        display: flex;
        flex-direction: column;
        gap: 10px;
        margin-top: 1%;
    }

    .file-item {
        display: flex;
        justify-content: space-between;
        align-items: center;
        background-color: #3a3a3a;
        border-radius: 5px;
        padding: 10px;
        color: white;
    }

    .file-item img.thumbnail {
        max-width: 50px;
        max-height: 50px;
        margin-right: 10px;
    }

    .file-item span {
        font-family: var(--code-font);
        flex-grow: 1;
        margin-right: 10px;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
    }

    .file-item button {
        background-color: var(--main-color);
        color: white;
        border: none;
        padding: 5px 10px;
        border-radius: 5px;
        cursor: pointer;
        font-family: var(--main-font);
        font-size: var(--font-size);
    }
</style>
