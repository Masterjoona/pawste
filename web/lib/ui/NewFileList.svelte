<script lang="ts">
    import "../../styles/file.css";
    export let files: File[] = [];
    export let imageSources = [];
    export let removeFile: (filename: string) => void;
    export let truncateFilename: (filename: string) => string;
    export let prettifyFileSize: (size: number) => string;
</script>

<div class="file-list">
    {#each files as file, index}
        <div class="file-item">
            {#if file.type.startsWith("image/")}
                <img
                    src={imageSources[index]}
                    alt={file.name}
                    class="thumbnail" />
            {/if}
            <span
                >{truncateFilename(file.name)} - {prettifyFileSize(
                    file.size,
                )}</span>
            <button on:click={() => removeFile(file.name)}>Remove</button>
        </div>
    {/each}
</div>
