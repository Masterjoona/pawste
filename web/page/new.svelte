<script lang="ts">
    import FileList from "../lib/ui/FileList.svelte";
    import {
        failToast,
        prettifyFileSize,
        handleAttachFiles,
    } from "../lib/utils.js";
    import "../styles/buttons.css";
    import "../styles/file.css";
    import "../styles/paste.css";

    export let fileUpload: boolean;
    export let maxFileSize: number;
    export let MaxEncryptionSize: number;
    export let maxContentLength: number;

    let selectedExpiration = "1w";
    let selectedBurnAfter = "0";
    let selectedSyntax = "none";
    let selectedPrivacy = "public";
    let content = "";
    let password = "";

    let attachedFiles: File[] = [];
    let imageSources: (string | null)[] = [];

    const setNewFiles = (files: File[]) => {
        attachedFiles = files;
    };

    const setImageSources = (sources: (string | null)[]) => {
        imageSources = sources;
    };

    function onFileAttach(event: any) {
        handleAttachFiles(
            event,
            attachedFiles,
            setNewFiles,
            imageSources,
            setImageSources,
        );
    }

    function handlePasswordChange(event: any) {
        password = event.target.value;
    }

    async function handleSave() {
        if (content === "" && attachedFiles.length === 0) {
            failToast("You must provide content or attach files!");
            return;
        }

        if (maxContentLength < content.length) {
            failToast(
                `Content is too long! Max length is ${maxContentLength} characters.`,
            );
            return;
        }

        const encrypted =
            selectedPrivacy === "private" || selectedPrivacy === "secret";

        if (encrypted && password === "") {
            failToast("You must provide a password for encrypted pastes!");
            return;
        }

        if (attachedFiles.length > 0) {
            for (let file of attachedFiles) {
                if (file.size > (encrypted ? MaxEncryptionSize : maxFileSize)) {
                    failToast(
                        `File ${file.name} is too large! Max size is ${
                            encrypted
                                ? prettifyFileSize(MaxEncryptionSize)
                                : prettifyFileSize(maxFileSize)
                        } bytes.`,
                    );
                    return;
                }
            }
        }

        const formData = new FormData();
        formData.append("expire", selectedExpiration);
        formData.append("burnafter", selectedBurnAfter);
        formData.append("syntax", selectedSyntax);
        formData.append("privacy", selectedPrivacy);
        formData.append("content", content);
        formData.append("password", password);
        console.log(attachedFiles);
        for (let file of attachedFiles) {
            formData.append("files[]", file);
        }

        const response = await (
            await fetch("/p", {
                method: "POST",
                body: formData,
            })
        ).json();

        if (response?.error) {
            failToast(response?.error);
        } else {
            setTimeout(() => {
                window.location.href = `/p/${response.pasteName}`;
            }, 500);
        }
    }

    function removeFile(filename: string) {
        attachedFiles = attachedFiles.filter((file) => file.name !== filename);
    }

    function handlePrivacyChange(event: any) {
        selectedPrivacy = event.target.value;
        const passwordField = document.getElementById("password-field");
        if (
            selectedPrivacy === "private" ||
            selectedPrivacy === "secret" ||
            selectedPrivacy === "readonly"
        ) {
            passwordField.style.display = "block";
        } else {
            passwordField.style.display = "none";
        }
    }
</script>


<svelte:head>
    <title>pawst.eu -- new paste</title>
</svelte:head>
<div id="container">
    <div class="card">
        <div class="options">
            <div>
                <label for="expiration">Expiration:</label>
                <select id="expiration" bind:value={selectedExpiration}>
                    <option value="never">Never</option>
                    <option value="1h">1 Hour</option>
                    <option value="6h">6 Hours</option>
                    <option value="1d">1 Day</option>
                    <option value="3d">3 Days</option>
                    <option value="1w">1 week</option>
                </select>
            </div>
            <div>
                <label for="burn-after">Burn After:</label>
                <select id="burn-after" bind:value={selectedBurnAfter}>
                    <option value="0">Never</option>
                    <option value="1">1 View</option>
                    <option value="10">10 Views</option>
                    <option value="100">100 Views</option>
                    <option value="1000">1000 Views</option>
                </select>
            </div>
            <div>
                <label for="syntax">Syntax:</label>
                <select id="syntax" bind:value={selectedSyntax}>
                    <option value="none">None</option>
                </select>
            </div>
            <div>
                <label for="privacy">Privacy:</label>
                <select
                    id="privacy"
                    bind:value={selectedPrivacy}
                    on:change={handlePrivacyChange}>
                    <option value="public">Public</option>
                    <option value="unlisted">Unlisted</option>
                    <option value="readonly">Read-only</option>
                    <option value="private">Private</option>
                    <option value="secret">Secret</option>
                </select>
            </div>
            <div id="password-field">
                <label for="password">Password:</label>
                <input
                    type="password"
                    id="password"
                    on:input={handlePasswordChange} />
            </div>
        </div>

        <textarea placeholder="Pawste away" bind:value={content}></textarea>
        <div class="buttons">
            {#if fileUpload}
                <input
                    type="file"
                    multiple
                    on:change={onFileAttach}
                    style="display: none;"
                    id="file-input" />
                <button
                    on:click={() =>
                        document.getElementById("file-input").click()}
                    >Attach Files</button>
            {/if}
            <button on:click={handleSave}>Save</button>
        </div>
        {#if fileUpload}<FileList
                files={attachedFiles}
                {imageSources}
                {removeFile} />
        {/if}
    </div>
</div>

<style>
    .options {
        display: flex;
        justify-content: space-evenly;
        flex-direction: column;
        gap: 10px;
        margin-bottom: 10px;
    }

    .options div {
        display: flex;
        align-items: center;
    }

    label {
        font-size: var(--font-size);
        color: white;
    }

    select {
        background-color: #1e1e1e;
        color: white;
        border: 1px solid #444;
        padding: 5px;
        border-radius: 5px;
        font-family: var(--main-font);
        font-size: var(--font-size);
    }

    #password-field {
        margin-top: -0.1%;
        display: none;
    }

    #password {
        width: 40%;
        padding: 5px;
        border-radius: 5px;
        border: 1px solid #444;
        background-color: #1e1e1e;
        color: white;
        font-family: var(--main-font);
        font-size: var(--font-size);
    }

    @media (min-width: 600px) {
        .options {
            flex-direction: row;
            flex-wrap: wrap;
            gap: 20px;
        }

        textarea {
            height: 400px;
        }
    }
</style>
