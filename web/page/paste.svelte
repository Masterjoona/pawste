<script>
    import { toast } from "@zerodevx/svelte-toast";
    import { copy } from "svelte-copy";
    import {
        truncateFilename,
        viewFile,
        timeDifference,
        prettifyFileSize,
    } from "../lib/utils.js";
    import "../styles/paste.css";
    import "../styles/file.css";
    import "../styles/buttons.css";

    export let paste;
    export let files;
    export let password;
    const successToast = (msg) => {
        toast.push(msg, {
            theme: {
                "--toastColor": "mintcream",
                "--toastBackground": "rgba(72,187,120,0.9)",
                "--toastBarBackground": "#2F855A",
            },
        });
    };
    async function deletePaste() {
        const resp = await fetch(`/p/${paste.PasteName}`, {
            method: "DELETE",
            body: JSON.stringify({ password }),
        });
        if (!resp.ok) {
            toast.push("Failed to delete!", {
                theme: {
                    "--toastColor": "mintcream",
                    "--toastBackground": "rgba(255,0,0,0.9)",
                    "--toastBarBackground": "red",
                },
            });
        } else {
            location.href = "/";
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
            <button on:click={deletePaste}>Delete</button>
        </div>
        <div class="file-list">
            {#each files as file}
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
