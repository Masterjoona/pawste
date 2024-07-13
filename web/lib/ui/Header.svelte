<script lang="ts">
    export let PublicList: boolean;
</script>

<center>
    <h1>pawste</h1>
    <h3>**currently in development expect bugs and data being wiped**</h3>
    <a href="/">new</a>
    {#if PublicList}
        <a href="/list">list</a>
    {/if}
    <a href="/about">about</a>
    <a href="https://github.com/sponsors/Masterjoona/">Sponsor!</a>
</center>

<style>
    h1 {
        font-family: var(--main-font), cursive;
        color: var(--main-color);
    }
    a {
        font-family: var(--code-font);
    }
</style>
