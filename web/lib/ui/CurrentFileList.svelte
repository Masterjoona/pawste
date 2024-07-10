<script lang="ts">
    import "../../styles/file.css";
    import { File } from "../types";
    import { prettifyFileSize, truncateFilename, viewFile } from "../utils";

    export let files: File[];
    export let pasteName: string;
    export let removeFile: (filename: string) => void = null;
</script>

<div class="file-list">
    {#each files as file}
        <div class="file-item">
            <span
                >{truncateFilename(file.Name)} - {prettifyFileSize(
                    file.Size,
                )}</span>
            {#if removeFile}
                <button on:click={() => removeFile(file.Name)}>Remove</button>
            {/if}
            <button on:click={() => viewFile(pasteName, file.Name)}
                >View</button>
        </div>
    {/each}
</div>
