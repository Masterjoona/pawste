<script lang="ts">
    import { Paste } from "../lib/types";
    import PasteList from "../lib/ui/PasteList.svelte";

    export let pastes: Paste[] = [];
    const headers = ["Name", "Expire", "View"];
</script>

<svelte:head>
    <title>pawste -- paste list</title>
</svelte:head>

<div id="spacemaker"></div>
<div id="container">
    <PasteList {pastes} tableHeaders={headers} />
</div>

<style>
    #spacemaker {
        margin-top: 50px;
    }
    #container {
        height: 100%;
        width: 100%;
        display: flex;
        flex-direction: column;
        justify-content: flex-start;
        align-items: center;
        font-family: var(--main-font);
        overflow: auto;
    }
</style>
