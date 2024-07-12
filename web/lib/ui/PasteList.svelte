<script lang="ts">
    import { successToast, timeDifference, deletePaste } from "../utils";
    import { Paste } from "../types";
    export let pastes: Paste[];
    export let tableHeaders: string[];
    export let password: string = "";
</script>

<table id="pastes-table">
    <tr>
        {#each tableHeaders as header}
            <th>{header}</th>
        {/each}
    </tr>
    {#if pastes === undefined || pastes.length === 0}
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
                {#if Privacy}
                    <td
                        >{timeDifference(ReadLast)}
                        <i class="fa-solid fa-clock"></i></td>
                    <td>{Privacy} <i class="fa-solid fa-lock"></i></td>
                    <td>{ReadCount} <i class="fa-solid fa-eye"></i></td>
                {/if}
                <td>
                    {#if UrlRedirect === 1}
                        <a href="/u/{PasteName}">Go to URL</a>
                    {:else}
                        <a href="/p/{PasteName}">View</a>
                    {/if}
                </td>
                {#if password}
                    <td>
                        <button
                            on:click={() =>
                                (window.location.href = "/e/" + PasteName)}
                            >Edit</button>
                        <button
                            on:click={() =>
                                deletePaste(PasteName, password, () =>
                                    successToast(`Deleted ${PasteName}!`),
                                )}>Delete</button>
                    </td>
                {/if}
            </tr>
        {/each}
    {/if}
</table>

<style>
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
