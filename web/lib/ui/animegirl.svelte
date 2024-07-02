<script>
    import { onMount } from "svelte";
    let src = "";
    let girl = "";

    const loadImage = async () => {
        const girls = [
            "faye",
            "asuka",
            "suzume",
            //"bocchi", too wide
            "mai",
            "marin",
            "nanami",
        ];
        girl = girls[~~(girls.length * Math.random())];
        const image = await import(`../../assets/${girl}.png`);
        src = image.default;
    };

    onMount(() => {
        loadImage();
    });
</script>

<img class="anime-girl" {src} alt={girl} />

<style>
    .anime-girl {
        position: fixed;
        max-height: 30%;
        bottom: 1%;
        right: -2%;
        margin: 10px;
    }
</style>
