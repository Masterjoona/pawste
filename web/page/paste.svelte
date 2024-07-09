<script>
    import { copy } from "svelte-copy";

    import {
        failToast,
        prettifyFileSize,
        successToast,
        timeDifference,
        truncateFilename,
        viewFile,
        deletePaste,
    } from "../lib/utils.js";

    import "../styles/buttons.css";
    import "../styles/file.css";
    import "../styles/password.css";
    import "../styles/paste.css";

    export let isEncrypted;
    export let paste;
    let inputPassword = "";
    let showContent = !isEncrypted;
    const triedPasswds = [];

    async function fetchPaste() {
        if (inputPassword === "") {
            failToast("Password cannot be empty!");
            return;
        } else if (triedPasswds.includes(inputPassword)) {
            failToast("Password already tried!");
            return;
        } else {
            const resp = await fetch(
                `/p/${window?.location?.pathname?.split("/").pop()}/json`,
                {
                    method: "GET",
                    headers: {
                        password: inputPassword,
                    },
                },
            );
            if (resp.ok) {
                const data = await resp.json();
                paste = { ...data };
                showContent = true;
            } else {
                failToast("Incorrect password!");
                triedPasswds.push(inputPassword);
            }
        }
    }
</script>

{#if isEncrypted && !showContent}
    <div class="overlay"></div>
    <div class="password-prompt">
        <label for="password">Enter password:</label>
        <input type="password" id="password" bind:value={inputPassword} />
        <button on:click={fetchPaste}>Submit</button>
    </div>
{/if}

<div id="container" class:blur={isEncrypted && !showContent}>
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
        <textarea readonly>{paste.Content}</textarea>
        <div class="buttons">
            <button
                on:click={() =>
                    (window.location.href = "/e/" + paste.PasteName)}
                >Edit</button>
            <button
                use:copy={paste.Content}
                on:svelte-copy={() => {
                    successToast("Text copied!");
                }}>Copy Text</button>
            <button
                use:copy={window?.location?.href}
                on:svelte-copy={() => {
                    successToast("URL copied!");
                }}>Copy URL</button>
            <button
                on:click={() =>
                    deletePaste(paste.PasteName, () => (location.href = "/"))}
                >Delete</button>
        </div>
        <div class="file-list">
            {#each paste.Files as file}
                <div class="file-item">
                    <span
                        >{truncateFilename(file.Name)} - {prettifyFileSize(
                            file.Size,
                        )}</span>
                    <button
                        on:click={() => viewFile(paste.PasteName, file.Name)}
                        >View</button>
                </div>
            {/each}
        </div>
    </div>
</div>
