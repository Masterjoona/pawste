<script>
    export let pastes;
</script>

<div id="spacemaker"></div>
<div id="container">
    {#if pastes?.length > 0}
        <table id="pastes">
            <tr>
                <th>Name</th>
                <th>Views</th>
                <th><i class="fa-solid fa-maximize"></i></th>
            </tr>
            {#each pastes as { ID, PasteName, ReadCount }}
                <tr>
                    <td>{PasteName}</td>
                    <td>{ReadCount} <i class="fa-solid fa-eye"></i></td>
                    <td><a href="/p/{PasteName}">View</a></td>
                </tr>
            {/each}
        </table>
    {:else}
        <div id="nopastes">
            <p>No pastes to show :(</p>
            <p>Would you like to <a href="/">make one</a> instead?</p>
        </div>
    {/if}
</div>

<style>
    #spacemaker {
        margin-top: 50px;
    }
    #container {
        height: 100%;
        width: 100%;
        display: flex;
        flex-direction: column;
        justify-content: flex-start;
        align-items: center;
        font-family: var(--main-font);
        overflow: auto;
    }
    table {
        width: 66.6%;
    }
    th,
    td {
        border: 1px solid var(--main-color);
        text-align: center;
    }
    #nopastes {
        border: 1px solid var(--main-color);
        border-radius: 15px;
        padding: 10px;
    }
</style>
