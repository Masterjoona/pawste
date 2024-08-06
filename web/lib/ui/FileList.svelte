<script lang="ts">
    import "../../styles/file.css";
    import type { FileType } from "../types";
    import {
        prettifyFileSize,
        truncateFilename,
        viewFile,
        isFileDb as FileDatafromDb,
    } from "../utils";

    export let files: FileType[] = [];
    export let imageSources = [];
    export let pasteName: string = "";
    export let removeFile: (filename: string) => void = null;
</script>

<div class="file-list">
    {#each files as file, index}
        {@const dataFromDb = FileDatafromDb(file)}
        <div class="file-item">
            {#if !dataFromDb && file.type.startsWith("image/")}
                <img
                    src={imageSources[index]}
                    alt={file.name}
                    class="thumbnail" />
            {/if}
            <span>
                {truncateFilename(dataFromDb ? file.Name : file.name)}
                ({prettifyFileSize(dataFromDb ? file.Size : file.size)})
            </span>
            {#if removeFile}
                <button
                    on:click={() =>
                        removeFile(dataFromDb ? file.Name : file.name)}
                    >Remove</button>
            {/if}
            {#if dataFromDb}
                <button on:click={() => viewFile(pasteName, file.Name)}
                    >View</button>
            {/if}
        </div>
    {/each}
</div>
