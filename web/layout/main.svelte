<script lang="ts">
    import { SvelteToast } from "@zerodevx/svelte-toast";

    import Animegirl from "../lib/ui/Animegirl.svelte";
    import Footer from "../lib/ui/Footer.svelte";
    import Header from "../lib/ui/Header.svelte";

    import "@fontsource-variable/fira-code";
    import "@fontsource/short-stack";
    import "@fortawesome/fontawesome-free/css/all.css";
    import "../app.css";

    export let AnimeGirls: boolean;
    export let PublicList: boolean;
</script>

<div id="container">
    <div id="header"><Header {PublicList} /></div>

    <main>
        <SvelteToast />

        <slot />
    </main>

    <div>
        {#if AnimeGirls}
            <Animegirl />
        {/if}
        <Footer />
    </div>
</div>

<style>
    #container {
        width: 100%;
        height: 100%;
        display: flex;
        flex-direction: column;
        align-items: center;
    }
    main {
        width: 80vw;
        height: 80vh;
    }
</style>
