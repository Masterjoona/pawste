<script lang="ts">
    import { url } from "golte/stores";
    import type { FileDb } from "../lib/types";

    import Password from "../lib/ui/Password.svelte";

    let hideContent = true;
    let password: string;
    let fileMetaData: FileDb;
    let isVideo = false;
    let isImage = false;
    let fileBytes: ArrayBuffer;

    async function onSubmit(inputPassword: string) {
        const resp = await fetch($url + "/json", {
            headers: {
                "Content-Type": "application/json",
                password: inputPassword,
            },
        });
        if (!resp.ok) return;
        password = inputPassword;
        fileMetaData = await resp.json();
        hideContent = false;
        isVideo = fileMetaData.ContentType.startsWith("video");
        isImage = fileMetaData.ContentType.startsWith("image");
        if (isVideo || isImage) {
            const resp = await fetch($url, {
                headers: {
                    password,
                },
            });
            if (!resp.ok) return;
            fileBytes = await resp.arrayBuffer();
        }
    }

    async function downloadFile() {
        const resp = await fetch($url, {
            headers: {
                password,
            },
        });
        if (!resp.ok) return;
        const blob = await resp.blob();
        const url = URL.createObjectURL(blob);
        const a = document.createElement("a");
        a.href = url;
        a.download = fileMetaData.Name;
        a.click();
        URL.revokeObjectURL(url);
    }
</script>

<div class="auth-file-container">
    {#if hideContent}
        <Password {onSubmit} question="Enter password to download/view file:" />
    {:else if fileMetaData}
        <h2>{fileMetaData.Name}</h2>
        <p>Size: {fileMetaData.Size}</p>
        <p>Content Type: {fileMetaData.ContentType}</p>
        {#if isVideo}
            <video controls>
                <source
                    type={fileMetaData.ContentType}
                    src={URL.createObjectURL(
                        new Blob([fileBytes], {
                            type: fileMetaData.ContentType,
                        }),
                    )} />
                <track kind="captions" />
            </video>
        {:else if isImage}
            <img
                src={URL.createObjectURL(
                    new Blob([fileBytes], { type: fileMetaData.ContentType }),
                )}
                alt={fileMetaData.Name} />
        {:else}
            <p>File type not supported for preview</p>
            <button on:click={downloadFile}>Download</button>
        {/if}
    {/if}
</div>
