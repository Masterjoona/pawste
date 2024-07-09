<main>
    <h1>Authentication</h1>
    <form
        id="auth-form"
        method="POST"
        action={"/"}
        enctype="multipart/form-data">
        <label for="password">Password:</label>
        <input
            id="password-field"
            placeholder="Password"
            name="password"
            type="password"
            autocomplete="off" />
        <button type="submit">Submit</button>
    </form>
    <script>
        document
            .getElementById("auth-form")
            .addEventListener("submit", function (event) {
                this.action = window.location.href.slice(0, -5);
            });
    </script>
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
