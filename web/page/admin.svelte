<script lang="ts">
    import FileSizeInput from "../lib/ui/FileSizeInput.svelte";
    import Password from "../lib/ui/Password.svelte";
    import StringInput from "../lib/ui/StringInput.svelte";
    import Switch from "../lib/ui/Switch.svelte";

    import { Config, Paste } from "../lib/types";
    import PasteList from "../lib/ui/PasteList.svelte";
    import { failToast, prettifyFileSize, successToast } from "../lib/utils";
    import "../styles/buttons.css";

    export let config: Config;
    export let pastes: Paste[];
    let adminPassword = "";

    async function deletePaste(pasteName: string) {
        const resp = await fetch(`/p/${pasteName}`, {
            method: "DELETE",
            body: JSON.stringify({ password: adminPassword }), // cursed but whatever
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

    function updateConfigValue(key: string, value: string | boolean) {
        config = { ...config, [key]: !value };
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
        <PasteList
            {pastes}
            tableHeaders={[
                "Name",
                "Expire",
                "Read Last",
                "Privacy",
                "Views",
                "",
                "Actions",
            ]}
            {deletePaste} />
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
                                onChange={() =>
                                    updateConfigValue(key, value)} />
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
                                    updateConfigValue(key, newValue)} />
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

    #pastes-section {
        flex: 1 1 100%;
        margin-top: 20px;
        display: flex;
        flex-direction: column;
        align-items: center;
    }

    #env-vars-section {
        margin-top: 20px;
        margin-bottom: 10%;
    }

    #env-vars-table,
    #info-table {
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
