<svelte:head>
    <title>pawste -- about</title>
</svelte:head>

<div>
    <h2>About us ♡</h2>
    <p>
        This project is open source and you can find the source code on <a
            href="https://github.com/Masterjoona/pawste "
            target="_blank">GitHub</a
        >! The backend is done by
        <a href="https://masterjoona.dev" target="_blank">Joona</a>
        and is written in Go and the basis for the frontend is done in Svelte by a friend.
    </p>
    <p>
        Send us an email at <a href="mailto:team@pawst.eu">team@pawst.eu</a>!
    </p>
</div>

<style>
    div {
        height: 100%;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
    }
    h2,
    p {
        text-align: center;
        width: 76.6%;
        align-self: center;
        line-height: 1.9em;
    }

    p {
        font-size: 1.2em;
    }
    a {
        color: rgb(169, 95, 120);
    }
</style>
