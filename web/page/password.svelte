<script>
    import { toast } from "@zerodevx/svelte-toast";

    let password = "";

    const handleSubmit = async (event) => {
        event.preventDefault();
        if (password === "correct_password") {
            console.log("Password is correct, redirecting to the secret page.");
        } else {
            toast.push("Invalid password, please try again.", {
                theme: {
                    "--toastBackground": "red",
                    "--toastColor": "white",
                    "--toastBarBackground": "red",
                },
            });
        }
    };
</script>

<form on:submit={handleSubmit}>
    <label for="password">Password:</label>
    <input type="password" id="password" bind:value={password} required />
    <button type="submit">Submit</button>
</form>

<style>
    form {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        height: 100%;
    }
    label {
        margin-bottom: 10px;
    }
    input {
        margin-bottom: 10px;
    }
</style>
