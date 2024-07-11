<script lang="ts">
    import "../../styles/file.css";
    import { FileType } from "../types";
    import {
        prettifyFileSize,
        truncateFilename,
        viewFile,
        isFileDb as isFromDb,
    } from "../utils";

    export let files: FileType[] = [];
    export let imageSources = [];
    export let pasteName: string = "";
    export let removeFile: (filename: string) => void = null;
</script>

<div class="file-list">
    {#each files as file, index}
        <div class="file-item">
            {#if !isFromDb(file) && file.type.startsWith("image/")}
                <img
                    src={imageSources[index]}
                    alt={file.name}
                    class="thumbnail" />
            {/if}
            {#if isFromDb(file)}
                <span
                    >{truncateFilename(file.Name)} - {prettifyFileSize(
                        file.Size,
                    )}</span>
                {#if removeFile}
                    <button on:click={() => removeFile(file.Name)}
                        >Remove</button>
                {/if}
                <button on:click={() => viewFile(pasteName, file.Name)}
                    >View</button>
            {:else}
                <span
                    >{truncateFilename(file.name)} - {prettifyFileSize(
                        file.size,
                    )}</span>
                {#if removeFile}
                    <button on:click={() => removeFile(file.name)}
                        >Remove</button>
                {/if}
                <button on:click={() => viewFile(pasteName, file.name)}
                    >View</button>
            {/if}
        </div>
    {/each}
</div>
