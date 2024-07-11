<script lang="ts">
    import FileList from "../lib/ui/FileList.svelte";
    import Password from "../lib/ui/Password.svelte";
    import Properties from "../lib/ui/Properties.svelte";

    import { Paste } from "../lib/types";
    import { failToast } from "../lib/utils";
    import "../styles/buttons.css";
    import "../styles/file.css";
    import "../styles/paste.css";
    export let paste: Paste;
    export let isEncrypted: boolean;

    let pastePassword: string;
    let showContent = !isEncrypted;
    let removedFiles: string[] = [];
    let newFiles: File[] = [];
    let imageSources = [];
    let newContent = paste.Content;
    const question = "Password needed to edit";

    function handleAttachFiles(event: any) {
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

    function removeOldFile(filename: string) {
        removedFiles = [...removedFiles, filename];
        paste.Files = paste.Files.filter((file) => file.Name !== filename);
    }
    function removeNewFile(filename: string) {
        newFiles = newFiles.filter((file) => file.name !== filename);
    }

    function filenamesConflict() {
        for (let file of newFiles) {
            if (
                paste.Files.some((f) => f.Name === file.name) &&
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
        return !paste.Files.length && !paste.Content && noNewFiles;
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

        const resp = await fetch(`/p/${paste.PasteName}`, {
            method: "PATCH",
            body: formData,
        });

        if (!resp.ok) {
            failToast("Failed to save!");
        } else {
            location.href = `/p/${paste.PasteName}`;
        }
    }

    async function fetchPaste(password: string) {
        const resp = await fetch(
            `/p/${window?.location?.pathname?.split("/").pop()}/json`,
            {
                method: "GET",
                headers: {
                    password: password,
                },
            },
        );
        if (resp.ok) {
            paste = await resp.json();
            newContent = paste.Content;
            showContent = true;
        } else {
            failToast("Something went wrong!");
        }
    }
</script>

{#if isEncrypted && !showContent}
    <Password {question} onSubmit={fetchPaste} />
{/if}

<div id="container">
    <div class="card">
        <Properties {paste} />
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
        <p>Current Files:</p>
        {#if paste.Files}
            <FileList
                files={paste.Files}
                pasteName={paste.PasteName}
                removeFile={removeOldFile} />
        {/if}
        <div class="file-list">
            <p>New Files:</p>
            <FileList
                files={newFiles}
                {imageSources}
                removeFile={removeNewFile} />
        </div>
    </div>
</div>
