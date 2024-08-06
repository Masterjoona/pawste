<script lang="ts">
    import { url } from "golte/stores";
    import { FileDb } from "../lib/types";

    import Password from "../lib/ui/Password.svelte";

    let password: string;
    let hideContent = true;
    let fileData: FileDb;

    async function onSubmit(inputPassword: string) {
        const resp = await fetch($url + "/json", {
            headers: {
                "Content-Type": "application/json",
                password: inputPassword,
            },
        });
        if (!resp.ok) return;
        password = inputPassword;
        fileData = await resp.json();
        hideContent = false;
    }
</script>

<div class="auth-file-container">
    {#if hideContent}
        <Password {onSubmit} question="Enter password to download/view file:" />
    {:else if fileData}
        <h2>{fileData.Name}</h2>
        <p>Size: {fileData.Size}</p>
    {/if}
</div>
