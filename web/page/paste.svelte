<script>
    import { toast } from "@zerodevx/svelte-toast";
    import { copy } from "svelte-copy";
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

    function truncateFilename(filename, maxLength = 30) {
        const extIndex = filename.lastIndexOf(".");
        const name = filename.substring(0, extIndex);
        const ext = filename.substring(extIndex);

        if (name.length + ext.length <= maxLength) {
            return filename;
        }

        const charsToShow = maxLength - ext.length - 3;
        const startChars = Math.ceil(charsToShow / 2);
        const endChars = Math.floor(charsToShow / 2);

        return (
            name.substring(0, startChars) +
            "..." +
            name.substring(name.length - endChars) +
            ext
        );
    }

    function viewFile(filename) {
        window.open("/p/" + paste.PasteName + "/f/" + filename);
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
                    <button on:click={() => viewFile(file.Name)}>View</button>
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
</style>
