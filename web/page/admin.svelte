<script lang="ts">
    import Password from "../lib/ui/Password.svelte";
    import StringInput from "../lib/ui/StringInput.svelte";
    import Switch from "../lib/ui/Switch.svelte";
    import FileSizeInput from "../lib/ui/FileSizeInput.svelte";

    import { Config, Paste } from "../lib/types";
    import {
        failToast,
        prettifyFileSize,
        successToast,
        timeDifference,
    } from "../lib/utils";
    import "../styles/buttons.css";

    export let config: Config;
    export let pastes: Paste[];
    let adminPassword = "";

    async function deletePaste(pasteName: string) {
        const resp = await fetch(`/p/${pasteName}`, {
            method: "DELETE",
            body: JSON.stringify({ adminPassword }), // cursed but whatever
        });
        if (!resp.ok) {
            failToast("Failed to delete paste!");
        } else {
            successToast(`Deleted ${pasteName}!`);
        }
    }

    async function fetchConfigAndPastes(password: string) {
        const resp = await fetch(location.pathname + "/json", {
            method: "GET",
            headers: { password },
        });
        if (resp.ok) {
            const data = await resp.json();

            delete data.config.AdminPassword;
            delete data.config.Salt;
            delete data.config.Port;

            config = data.config;
            pastes = data.pastes;
            adminPassword = password;
        } else {
            failToast("Incorrect password!");
        }
    }

    function toggleConfigBool(key: string, value: boolean) {
        config = { ...config, [key]: !value };
    }

    function updateConfigString(key: string, newValue: string) {
        config = { ...config, [key]: newValue };
    }

    function updateConfigNumber(key: string, event: Event) {
        const newValue = parseInt((event.target as HTMLInputElement).value);
        config = { ...config, [key]: newValue };
    }
</script>

{#if !adminPassword}
    <Password question={"Admin password"} onSubmit={fetchConfigAndPastes} />
{/if}

<div id="admin-container">
    <h2>pawste admin</h2>

    <div id="links-section">
        <h3>Links</h3>
        <ul>
            <li>
                <a href="https://github.com/Masterjoona/pawste">Source Code</a>
            </li>
            <li>
                <a href="https://github.com/Masterjoona/pawste/issues"
                    >Feedback</a>
            </li>
            <li>
                <a href="https://github.com/sponsors/Masterjoona/">Sponsor</a>
            </li>
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
            {#if pastes === null || pastes.length === 0}
                <tr>
                    <td colspan="7">No pastes</td>
                </tr>
            {:else}
                {#each pastes as { PasteName, ReadCount, ReadLast, Privacy, Expire, UrlRedirect }}
                    <tr>
                        <td>{PasteName}</td>
                        <td
                            >{timeDifference(Expire)}
                            <i class="fa-solid fa-clock"></i></td>
                        <td
                            >{timeDifference(ReadLast)}
                            <i class="fa-solid fa-clock"></i></td>
                        <td>{Privacy} <i class="fa-solid fa-lock"></i></td>
                        <td>{ReadCount} <i class="fa-solid fa-eye"></i></td>
                        <td>
                            {#if UrlRedirect === 1}
                                <a href="/u/{PasteName}">Go to URL</a>
                            {:else}
                                <a href="/p/{PasteName}">View</a>
                            {/if}
                        </td>
                        <td>
                            <button
                                on:click={() =>
                                    (window.location.href = "/e/" + PasteName)}
                                >Edit</button>
                            <button on:click={() => deletePaste(PasteName)}
                                >Delete</button>
                        </td>
                    </tr>
                {/each}
            {/if}
        </table>
    </div>

    <div id="env-vars-section">
        <h3>Environmental Variables</h3>
        <table id="env-vars-table">
            <tr>
                <th>Argument</th>
                <th>Value</th>
                <th>Toggle</th>
            </tr>
            {#each Object.entries(config) as [key, value]}
                {#if typeof value === "boolean"}
                    <tr>
                        <td>{key}</td>
                        <td>{value}</td>
                        <td>
                            <Switch
                                checked={value}
                                onChange={() => toggleConfigBool(key, value)} />
                        </td>
                    </tr>
                {:else if typeof value === "string"}
                    <tr>
                        <td>{key}</td>
                        <td>{value}</td>
                        <td>
                            <StringInput
                                {value}
                                onInput={(newValue) =>
                                    updateConfigString(key, newValue)} />
                        </td>
                    </tr>
                {:else if key === "MaxEncryptionSize" || key === "MaxFileSize"}
                    <tr>
                        <td>{key}</td>
                        <td>
                            {prettifyFileSize(value)}
                        </td>
                        <td
                            ><FileSizeInput
                                {value}
                                onInput={(newValue) =>
                                    updateConfigNumber(key, newValue)} /></td>
                    </tr>
                {:else if typeof value === "number"}
                    <tr>
                        <td>{key}</td>
                        <td>{value}</td>
                        <td>
                            <input
                                type="number"
                                {value}
                                on:input={(e) => updateConfigNumber(key, e)} />
                        </td>
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
        justify-content: center;
    }

    #admin-container h2 {
        width: 100%;
        text-align: center;
    }

    #links-section,
    #info-section {
        flex: 1 1 calc(50% - 20px);
    }

    #pastes-section {
        flex: 1 1 100%;
        margin-top: 20px;
        display: flex;
        flex-direction: column;
        align-items: center;
    }

    #env-vars-section {
        flex: 1 1 100%;
        margin-top: 20px;
    }

    #env-vars-table,
    #info-table,
    #pastes-table {
        width: 80%;
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
