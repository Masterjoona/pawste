<script lang="ts">
    import { copy } from "svelte-copy";

    import Currentfiles from "../lib/ui/CurrentFileList.svelte";
    import Password from "../lib/ui/Password.svelte";
    import Properties from "../lib/ui/Properties.svelte";

    import { Paste } from "../lib/types";
    import { failToast, successToast } from "../lib/utils";

    import "../styles/buttons.css";
    import "../styles/file.css";
    import "../styles/password.css";
    import "../styles/paste.css";

    export let isEncrypted: boolean;
    export let paste: Paste;
    export let burnAfter: boolean;

    let password: string;
    let showContent = !isEncrypted;
    let question = "Enter password:";

    question = burnAfter
        ? (question += " (Will be burned after read)")
        : question;

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
            showContent = true;
        } else {
            failToast("Incorrect password!");
        }
    }

    async function deletePaste(pasteName: string) {
        if (paste.Privacy === "readonly") {
            showContent = false;
            isEncrypted = true;
            question = "Password needed to delete paste:";
        }
        const resp = await fetch(`/p/${pasteName}`, {
            method: "DELETE",
            body: JSON.stringify({ password }),
        });
        if (!resp.ok) {
            failToast("Failed to delete paste!");
        } else {
            location.href = "/";
        }
    }
</script>

{#if isEncrypted && !showContent}
    <Password {question} onSubmit={fetchPaste} />
{/if}

<div id="container">
    <div class="card">
        <Properties {paste} />
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
            <button on:click={() => deletePaste(paste.PasteName)}
                >Delete</button>
        </div>
        <Currentfiles files={paste.Files} pasteName={paste.PasteName} />
    </div>
</div>
