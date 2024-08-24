<script lang="ts">
    import "../../styles/password.css";
    import { failToast } from "../utils";

    export let question = "Enter password:";
    export let onSubmit: (password: string) => void;
    let inputPassword: string = "";
    const triedPasswds = [];

    function callSubmit() {
        if (inputPassword === "") {
            failToast("Password cannot be empty!");
            return;
        } else if (triedPasswds.includes(inputPassword)) {
            failToast("Password already tried!");
            return;
        } else {
            triedPasswds.push(inputPassword);
            onSubmit(inputPassword);
        }
    }
</script>

<div class="overlay"></div>
<div class="password-prompt">
    <label for="password">{question}</label>
    <input type="password" id="password" bind:value={inputPassword} />
    <button on:click={callSubmit} class="submit">Submit</button>
</div>
