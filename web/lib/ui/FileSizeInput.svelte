<script lang="ts">
    import { prettifyFileSize } from "../utils";
    export let value;
    export let onInput;

    let editableValue = value;

    function handleInput(event) {
        editableValue = event.target.value;
        onInput(parseInt(editableValue));
    }
</script>

<div>
    <input
        type="number"
        value={editableValue}
        on:input={handleInput}
        on:blur={() => (editableValue = value)} />
</div>
