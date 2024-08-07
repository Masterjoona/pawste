<script lang="ts">
    import FileList from "../lib/ui/FileList.svelte";
    import Password from "../lib/ui/Password.svelte";
    import Properties from "../lib/ui/Properties.svelte";

    import type { Paste } from "../lib/types";
    import { failToast, handleAttachFiles, savePaste } from "../lib/utils";

    import "../styles/buttons.css";
    import "../styles/file.css";
    import "../styles/paste.css";

    export let paste: Paste;
    export let needsAuth: boolean;

    let pastePassword: string;
    let hideContent = needsAuth;
    let removedFiles: string[] = [];
    let newFiles: File[] = [];
    let imageSources: (string | null)[] = [];
    let newContent = paste.Content;
    const question = "Password needed to edit";

    const addNewFile = (file: File) => {
        newFiles = [...newFiles, file];
    };
    const addImageSources = (source: string) => {
        imageSources = [...imageSources, source];
    };

    function onFileAttach(event: any) {
        handleAttachFiles(event, addNewFile, addImageSources);
    }

    function removeNewFile(filename: string) {
        const fileIndex = newFiles.findIndex((file) => file.name === filename);
        if (fileIndex !== -1) {
            newFiles = newFiles.filter((_, index) => index !== fileIndex);
            imageSources = imageSources.filter(
                (_, index) => index !== fileIndex,
            );
        }
    }

    function removeOldFile(filename: string) {
        removedFiles = [...removedFiles, filename];
        paste.Files = paste.Files.filter((file) => file.Name !== filename);
    }

    function filenamesConflict() {
        for (let file of newFiles) {
            if (
                paste.Files?.some((f) => f.Name === file.name) &&
                !removedFiles.includes(file.name)
            ) {
                return true;
            }
        }
        return false;
    }

    const noChanges = (noNewFiles: boolean) => {
        return (
            newContent === paste.Content && noNewFiles && !removedFiles.length
        );
    };

    const emptyPaste = (noNewFiles: boolean) => {
        return !paste.Files?.length && !paste.Content && noNewFiles;
    };

    async function handleSave() {
        const noNewFiles = !newFiles.length;
        if (emptyPaste(noNewFiles)) {
            failToast("You must provide content or attach files!");
            return;
        }
        if (noChanges(noNewFiles)) {
            failToast("No changes detected!");
            return;
        }

        if (filenamesConflict()) {
            failToast("Filenames conflict!");
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
        formData.append("password", pastePassword);

        const saveButton = document.getElementById("save-button");
        const url = `/p/${paste.PasteName}`;
        savePaste("PATCH", formData, url, saveButton);
    }

    async function fetchPaste(password: string) {
        const resp = await fetch(
            location.pathname.replace("/e/", "/p/") + "/json",
            {
                headers: {
                    password: password,
                },
            },
        );
        if (resp.ok) {
            paste = await resp.json();
            newContent = paste.Content;
            hideContent = false;
        } else {
            failToast("Wrong password!");
        }
    }
</script>

{#if needsAuth && hideContent}
    <Password {question} onSubmit={fetchPaste} />
{/if}

<div id="edit-container">
    <div class="card">
        <Properties {paste} />
        <textarea bind:value={newContent}></textarea>
        <div class="buttons">
            <input
                type="file"
                multiple
                on:change={onFileAttach}
                style="display: none;"
                id="file-input" />
            <button
                on:click={() => document.getElementById("file-input").click()}
                >Attach Files</button>
            <button id="save-button" on:click={handleSave}>Save</button>
        </div>
        <p>Current Files:</p>
        {#if paste.Files}
            <FileList
                files={paste.Files}
                pasteName={paste.PasteName}
                removeFile={removeOldFile} />
        {/if}
        <p>New Files:</p>
        <FileList files={newFiles} {imageSources} removeFile={removeNewFile} />
    </div>
</div>
