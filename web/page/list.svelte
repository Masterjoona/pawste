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
                <th>Expire</th>
                <th><i class="fa-solid fa-maximize"></i></th>
            </tr>
            {#each pastes as { PasteName, ReadCount, Expire, UrlRedirect }}
                <tr>
                    <td>{PasteName}</td>
                    <td>{Expire} <i class="fa-solid fa-clock"></i></td>
                    <td>{ReadCount} <i class="fa-solid fa-eye"></i></td>
                    <td>
                        {#if UrlRedirect === 1}
                            <a href="/u/{PasteName}">Go to URL</a>
                        {:else}
                            <a href="/p/{PasteName}">View</a>
                        {/if}
                    </td>
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
