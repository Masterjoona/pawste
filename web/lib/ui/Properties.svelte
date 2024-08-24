<script lang="ts">
    import type { Paste } from "../types";
    import { timeDifference } from "../utils";
    export let paste: Paste;
</script>

<div class="properties">
    <h>{paste.PasteName}</h>
    <div class="spacer"></div>
    <div class="icon-container">
        <p>{paste.ReadCount} <i class="fa-solid fa-eye"></i></p>
        <p>
            {timeDifference(paste.ReadLast, "Read last")}
        </p>
        {#if paste.BurnAfter !== 0}
            <p>
                Burn after {paste.BurnAfter} views
            </p>
        {/if}
        <p>
            {paste.Content.length}
            <i class="fa-solid fa-file-lines"></i>
        </p>
        <p>
            {timeDifference(paste.Expire, "Expires in")}
            <i class="fa-solid fa-clock"></i>
        </p>
    </div>
</div>
