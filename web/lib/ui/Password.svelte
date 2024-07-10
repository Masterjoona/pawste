<script>
    import { failToast } from "../utils";
    import "../../styles/password.css";
    export let question = "Enter password:";
    export let onSubmit;
    let inputPassword;

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
