<script>
    import "../styles/paste.css";
    import "../styles/buttons.css";

    export let isEncrypted;
    export let password;
    let onYes = () => {
        if (!isEncrypted) {
            location.href += "?read=true";
        } else {
            const form = document.createElement("form");
            form.method = "POST";
            form.action = location.pathname + "?read=true";
            form.style.display = "none";

            const passwordInput = document.createElement("input");
            passwordInput.type = "hidden";
            passwordInput.name = "password";
            passwordInput.value = password;
            form.appendChild(passwordInput);

            document.body.appendChild(form);
            form.submit();
        }
    };
    let onNo = () => {
        location.href = "/";
    };
</script>

<div class="container">
    <p>This paste will burn after this, do you want to see this?</p>
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
