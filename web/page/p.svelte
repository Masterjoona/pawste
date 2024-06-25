<script>
    export let pasteName;
    export let expire;
    export let content;
    export let readCount;
</script>

<div id="container">
    <div class="card">
        <div class="properties">
            <p>{pasteName} | commies (id: cummies)</p>
            <div class="spacer"></div>
            <div class="icon-container">
                <p>{readCount} <i class="fa-solid fa-eye"></i></p>
                <p>{content.length} <i class="fa-solid fa-file-lines"></i></p>
                <p>{expire} <i class="fa-solid fa-clock"></i></p>
            </div>
        </div>
        <textarea readonly>{content}</textarea>
        <div class="buttons">
            <button>Modify</button>
            <button>Copy Text</button>
            <button>Copy URL</button>
        </div>
    </div>
</div>

<style>
    :root {
        --font-size: 1.2em;
    }

    #container {
        height: 100%;
        width: 100%;
        display: flex;
        justify-content: center;
        align-items: center;
        font-family: var(--main-font);
    }

    .card {
        width: 90%;
        display: flex;
        justify-content: center;
        align-items: center;
        flex-direction: column;
        background-color: #2a2a2a;
        border-radius: 10px;
        padding: 16px;
    }

    .properties {
        width: 100%;
        display: flex;
        flex-direction: row;
        justify-content: space-evenly;
        margin-top: 0px;
    }
    .properties {
        width: 100%;
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        margin-top: 0px;
    }

    .spacer {
        flex-grow: 1;
    }

    .icon-container {
        display: flex;
        flex-direction: row;
        gap: 10px;
    }

    textarea {
        width: 99%;
        height: 200px;
        font-family: var(--code-font);
        background-color: #1b1b22;
        color: white;
        border: none;
        border-radius: 10px;
        padding: 10px;
        resize: vertical;
        margin-bottom: 10px;
    }

    .buttons {
        display: flex;
        width: 100%;
        justify-content: space-evenly;
    }

    button {
        background-color: var(--main-color);
        color: white;
        border: none;
        padding: 10px 10px;
        border-radius: 5px;
        cursor: pointer;
        font-family: var(--main-font);
        font-size: var(--font-size);
    }

    button:hover {
        background-color: var(--main-color-dark);
    }
</style>
