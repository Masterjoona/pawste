<script>
    export let error;
</script>

<div class="error-container">
    <p class="error-title">{error} &gt;w&lt;</p>
    <a class="error-details" href="/">Take me back to safety</a>
</div>


<style>
    .error-container {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        height: 100%;
    }

    .error-title {
        font-size: 2em;
        margin-bottom: 20px;
    }

    .error-details {
        font-size: 1.2em;
    }
</style>