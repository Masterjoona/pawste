<script>
    export let onYes = () => {
        location.href += "?read=true";
    };
    export let onNo = () => {
        location.href = "/";
    };
</script>

<div class="container">
    <p>This paste will burn after this, do you want to see this?</p>
    <div class="buttons">
        <button on:click={onNo}>No</button>
        <button on:click={onYes}>Yes</button>
    </div>
</div>

<style>
    .container {
        background-color: #2a2a2a;
        border-radius: 10px;
        padding: 16px;
        max-width: 25%;
        margin: auto;
        text-align: center;
        color: white;
        border: 1px solid #444;
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
    }

    .buttons {
        display: flex;
        justify-content: space-around;
        width: 100%;
        margin-top: 20px;
    }

    button {
        background-color: var(--main-color);
        color: white;
        border: none;
        padding: 10px 20px;
        border-radius: 5px;
        cursor: pointer;
        font-family: var(--main-font);
        font-size: var(--font-size);
    }

    button:hover {
        background-color: var(--main-color-dark);
    }

    @media (max-width: 600px) {
        .container {
            max-width: 100%;
        }
    }
</style>
