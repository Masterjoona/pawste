<script>
    import {
        truncateFilename,
        viewFile,
        timeDifference,
    } from "../lib/utils.js";
    import { toast } from "@zerodevx/svelte-toast";
    import "../styles/paste.css";
    import "../styles/file.css";

    export let paste;
    export let files;
    export let password;

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

    function filenamesConflict() {
        for (let file of newFiles) {
            if (
                files.some((f) => f.Name === file.name) &&
                !removedFiles.includes(file.name)
            ) {
                return true;
            }
        }
        return false;
    }

    async function handleSave() {
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

        if (filenamesConflict()) {
            toast.push("Filenames conflict!", {
                theme: {
                    "--toastColor": "mintcream",
                    "--toastBackground": "rgba(255,0,0,0.9)",
                    "--toastBarBackground": "red",
                },
            });
            return;
        }

        const formData = new FormData();
        formData.append("content", newContent);
        for (let file of newFiles) {
            formData.append("files[]", file);
        }
        removedFiles.forEach((file) => {
            formData.append("removed_files", file);
        });
        formData.append("password", password);

        const resp = await fetch(`/p/${paste.PasteName}`, {
            method: "PATCH",
            body: formData,
        });

        if (!resp.ok) {
            toast.push("Failed to save!", {
                theme: {
                    "--toastColor": "mintcream",
                    "--toastBackground": "rgba(255,0,0,0.9)",
                    "--toastBarBackground": "red",
                },
            });
        } else {
            location.href = `/p/${paste.PasteName}`;
        }
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
                <p>
                    {timeDifference(paste.Expire)}
                    <i class="fa-solid fa-clock"></i>
                </p>
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
                    <button
                        on:click={() => viewFile(paste.PasteName, file.Name)}
                        >View</button>
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
