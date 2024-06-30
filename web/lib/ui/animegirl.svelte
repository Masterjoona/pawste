<script>
    import { onMount } from "svelte";
    let src = "";

    const loadImage = async () => {
        const girls = ["faye", "asuka", "suzume"];
        const girl = girls[Math.floor(girls.length * Math.random())];
        const image = await import(`../../assets/${girl}.png`);
        src = image.default;
    };

    onMount(() => {
        loadImage();
    });
</script>

<img class="anime-girl" {src} alt="girl" />

<style>
    .anime-girl {
        position: fixed;
        max-height: 30%;
        bottom: 0;
        right: 0;
        margin: 10px;
    }
</style>
