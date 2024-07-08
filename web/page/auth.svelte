<script>
    let password = "";

    function handleSubmit(event) {
        event.preventDefault();
        const form = document.createElement("form");
        form.method = "POST";
        form.action = location.pathname.slice(0, -5);
        form.style.display = "none";

        const passwordInput = document.createElement("input");
        passwordInput.type = "hidden";
        passwordInput.name = "password";
        passwordInput.value = password;
        form.appendChild(passwordInput);

        document.body.appendChild(form);
        form.submit();
    }
</script>

<main>
    <h1>Authentication</h1>
    <form on:submit={handleSubmit}>
        <label for="password">Password:</label>
        <input type="password" id="password" bind:value={password} />
        <button type="submit">Submit</button>
    </form>
</main>

<style>
    main {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        height: 65vh;
    }

    h1 {
        margin-bottom: 1rem;
    }

    form {
        display: flex;
        flex-direction: column;
        align-items: center;
    }

    label {
        margin-bottom: 0.5rem;
    }

    input {
        padding: 0.5rem;
        margin-bottom: 1rem;
    }

    button {
        padding: 0.5rem 1rem;
    }
</style>
