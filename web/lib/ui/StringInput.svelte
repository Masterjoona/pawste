<script>
    export let value;
    export let onInput;

    function handleInput(event) {
        onInput(event.target.value);
    }
</script>

<input type="text" {value} on:input={handleInput} />
