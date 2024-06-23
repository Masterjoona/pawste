<div id="container">
    <div class="card">
        <div class="options">
            <div>
                <label for="expiration">Expiration:</label>
                <select id="expiration">
                    <option value="never">Never</option>
                    <option value="1h">1 Hour</option>
                    <option value="6h">6 Hours</option>
                    <option value="1d">1 Day</option>
                    <option value="7d">7 Days</option>
                </select>
            </div>
            <div>
                <label for="burn-after">Burn After:</label>
                <select id="burn-after">
                    <option value="never">Never</option>
                    <option value="1view">1 View</option>
                    <option value="10view">10 Views</option>
                </select>
            </div>
            <div>
                <label for="syntax">Syntax:</label>
                <select id="syntax">
                    <option value="none">None</option>
                </select>
            </div>
            <div>
                <label for="privacy">Privacy:</label>
                <select id="privacy">
                    <option value="public">Public</option>
                    <option value="unlisted">Unlisted</option>
                </select>
            </div>
        </div>
        <textarea></textarea>
        <div class="buttons">
            <button>Attach File</button>
            <button>Save</button>
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

    .options {
        display: flex;
        justify-content: space-evenly;
        flex-direction: column;
        gap: 10px;
        margin-bottom: 10px;
    }

    .options div {
        display: flex;
        align-items: center;
    }

    label {
        font-size: var(--font-size);
        color: white;
    }

    select {
        background-color: #1E1E1E;
        color: white;
        border: 1px solid #444;
        padding: 5px;
        border-radius: 5px;
        font-family: var(--main-font);
        font-size: var(--font-size);
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
        padding: 10px 20px;
        border-radius: 5px;
        cursor: pointer;
        font-family: var(--main-font);
        font-size: var(--font-size);
    }

    button:hover {
        background-color: var(--main-color-dark);
    }

    @media (min-width: 600px) {
        .options {
            flex-direction: row;
            flex-wrap: wrap;
            gap: 20px;
        }

        .options div {
            flex: 1;
        }

        textarea {
            height: 400px;
        }
    }
</style>