<script>
    import "@fontsource-variable/fira-code";
    import "@fontsource/short-stack";
    import "@fortawesome/fontawesome-free/css/all.css";
    import { SvelteToast } from "@zerodevx/svelte-toast";
    import "../app.css";
    import Animegirl from "../lib/ui/animegirl.svelte";
    import Footer from "../lib/ui/footer.svelte";
    import Header from "../lib/ui/header.svelte";
</script>

<div id="container">
    <div id="header"><Header /></div>

    <main>
        <SvelteToast />
        <Animegirl />
        <slot />
    </main>

    <div id="footer"><Footer /></div>
</div>

<style>
    #container {
        width: 100%;
        height: 100%;
        display: flex;
        flex-direction: column;
        align-items: center;
    }
    #footer {
        position: absolute;
        bottom: 0px;
        right: 0px;
    }
    main {
        width: 80vw;
        height: 80vh;
    }
</style>
