<script lang="ts">
    export let value: string | number;
    export let onInput: (value: string | number) => void;
    export let type: "string" | "number";

    const handleInput = (event: Event) => {
        onInput((event.target as HTMLInputElement).value);
    };
</script>

<input {type} {value} on:input={handleInput} />
