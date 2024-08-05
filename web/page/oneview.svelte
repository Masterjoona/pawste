<script lang="ts">
    import "../styles/buttons.css";
    import "../styles/paste.css";

    const onYes = () => {
        location.href += "?read=true";
    };
    const onNo = () => {
        location.href = "/";
    };
</script>

<div class="container">
    <p>This paste will be deleted after this, do you want to see this?</p>
    <div class="buttons">
        <button on:click={onNo}>No</button>
        <button on:click={onYes}>Yes</button>
    </div>
</div>

<style>
    @media (max-width: 600px) {
        .container {
            max-width: 100%;
        }
    }
</style>
