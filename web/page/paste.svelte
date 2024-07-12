<script lang="ts">
    import { copy } from "svelte-copy";

    import FileList from "../lib/ui/FileList.svelte";
    import Password from "../lib/ui/Password.svelte";
    import Properties from "../lib/ui/Properties.svelte";

    import { Paste } from "../lib/types";
    import { failToast, successToast } from "../lib/utils";

    import "../styles/buttons.css";
    import "../styles/file.css";
    import "../styles/password.css";
    import "../styles/paste.css";

    export let needsAuth: boolean;
    export let paste: Paste;
    export let burnAfter: boolean;

    let password: string;
    let hideContent = needsAuth && paste.Privacy !== "readonly";
    console.log(needsAuth, hideContent);
    let question = "Enter password:";

    question = burnAfter
        ? (question += " (Will be burned after read)")
        : question;

    async function fetchPaste(password: string) {
        const resp = await fetch(location.pathname + "/json", {
            method: "GET",
            headers: {
                password: password,
            },
        });
        if (resp.ok) {
            paste = await resp.json();
            hideContent = false;
        } else {
            failToast("Incorrect password!");
        }
    }

    async function deletePaste(password: string) {
        const resp = await fetch(`/p/${paste.PasteName}`, {
            method: "DELETE",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ password }),
        });
        if (!resp.ok) {
            failToast("Failed to delete paste!");
        } else {
            location.href = "/";
        }
    }

    let onSubmitFunc = fetchPaste;

    function handleDelete() {
        if (paste.Privacy === "readonly") {
            hideContent = true;
            question = "Password needed to delete paste:";
            onSubmitFunc = async (password) => {
                await deletePaste(password);
            };
        } else {
            console.log("Deleting paste");
            deletePaste(password);
        }
    }
</script>

{#if needsAuth && hideContent}
    <Password {question} onSubmit={onSubmitFunc} />
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
            <button on:click={handleDelete}>Delete</button>
        </div>
        <FileList
            files={paste.Files ? paste.Files : []}
            pasteName={paste.PasteName} />
    </div>
</div>
