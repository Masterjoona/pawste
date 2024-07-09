<script>
    import { toast } from "@zerodevx/svelte-toast";
    import { copy } from "svelte-copy";
    import {
        truncateFilename,
        viewFile,
        timeDifference,
    } from "../lib/utils.js";
    import "../styles/paste.css";
    import "../styles/file.css";

    export let paste;
    export let files;
    const successToast = (msg) => {
        toast.push(msg, {
            theme: {
                "--toastColor": "mintcream",
                "--toastBackground": "rgba(72,187,120,0.9)",
                "--toastBarBackground": "#2F855A",
            },
        });
    };
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
        </div>
        <div class="file-list">
            {#each files as file}
                <div class="file-item">
                    <span
                        >{truncateFilename(file.Name)} - {(
                            file.Size / 1024
                        ).toFixed(2)} KB</span>
                    <button
                        on:click={() => viewFile(paste.PasteName, file.Name)}
                        >View</button>
                </div>
            {/each}
        </div>
    </div>
</div>
