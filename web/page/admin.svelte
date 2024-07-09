<script>
    import { deletePaste, prettifyFileSize, successToast } from "../lib/utils";
    import "../styles/buttons.css";

    export let config;
    export let pastes;
</script>

<div id="admin-container">
    <h2>pawste admin</h2>

    <div id="links-section">
        <h3>Links</h3>
        <ul>
            <li><a href="/docs">Documentation and Help</a></li>
            <li><a href="/source">Source Code</a></li>
            <li><a href="/issues">Feedback</a></li>
            <li><a href="/donate">Donate and Sponsor</a></li>
        </ul>
    </div>

    <div id="info-section">
        <h3>Info</h3>
        <table id="info-table">
            <tr><td>Version</td><td>dev</td></tr>
            <tr><td>Status</td><td>nop</td></tr>
            <tr><td>Uploads</td><td>morbillion</td></tr>
            <tr><td>Update</td><td>nop</td></tr>
        </table>
    </div>

    <div id="pastes-section">
        <h3>Pastes</h3>
        <table id="pastes-table">
            <tr>
                <th>Name</th>
                <th>Expire</th>
                <th>Last Read</th>
                <th>Privacy</th>
                <th>Views</th>
                <th>Link</th>
                <th>Actions</th>
            </tr>
            {#each pastes as { PasteName, ReadCount, ReadLast, Privacy, Expire, UrlRedirect }}
                <tr>
                    <td>{PasteName}</td>
                    <td>{Expire} <i class="fa-solid fa-clock"></i></td>
                    <td>{ReadLast} <i class="fa-solid fa-clock"></i></td>
                    <td>{Privacy} <i class="fa-solid fa-lock"></i></td>
                    <td>{ReadCount} <i class="fa-solid fa-eye"></i></td>
                    <td>
                        {#if UrlRedirect === 1}
                            <a href="/u/{PasteName}">Go to URL</a>
                        {:else}
                            <a href="/p/{PasteName}">View</a>
                        {/if}
                    </td>
                    <td
                        ><button
                            on:click={() =>
                                (window.location.href = "/e/" + PasteName)}
                            >Edit</button
                        ><button
                            on:click={() =>
                                deletePaste(PasteName, () =>
                                    successToast(`Deleted ${PasteName}!`),
                                )}>Delete</button
                        ></td>
                </tr>
            {/each}
        </table>
    </div>

    <div id="env-vars-section">
        <h3>Environmental Variables</h3>
        <table id="env-vars-table">
            <tr>
                <th>Argument</th>
                <th>Value</th>
            </tr>
            {#each Object.entries(config) as [key, value]}
                {#if key === "MaxEncryptionSize" || key === "MaxFileSize"}
                    <tr>
                        <td>{key}</td>
                        <td>{prettifyFileSize(value)}</td>
                    </tr>
                {:else}
                    <tr>
                        <td>{key}</td>
                        <td>{value}</td>
                    </tr>
                {/if}
            {/each}
        </table>
    </div>
</div>

<style>
    #admin-container {
        display: flex;
        flex-wrap: wrap;
        gap: 20px;
        font-family: var(--main-font);
        margin-top: 20px;
    }

    #admin-container h2 {
        width: 100%;
        text-align: center;
    }

    #links-section,
    #info-section {
        flex: 1 1 calc(50% - 20px);
    }

    #env-vars-section {
        flex: 1 1 100%;
        margin-top: 20px;
    }

    #env-vars-table {
        font-family: var(--code-font);
    }

    table {
        width: 100%;
        border-collapse: collapse;
        margin-top: 10px;
    }

    table,
    th,
    td {
        border: 1px solid black;
        padding: 8px;
        text-align: left;
    }

    td {
        vertical-align: top;
    }
</style>
