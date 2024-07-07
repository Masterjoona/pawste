<script>
    import { toast } from "@zerodevx/svelte-toast";
    let password = "";

    async function handleSubmit(event) {
        event.preventDefault();
        const formData = new FormData();
        formData.append("password", password);

        try {
            const response = await fetch("/p/auth/thing", {
                method: "POST",
                body: formData,
            });

            if (response.ok) {
                console.log("Authentication successful");
                toast.push("Redirecting shortly!", {
                    theme: {
                        "--toastColor": "mintcream",
                        "--toastBackground": "rgba(72,187,120,0.9)",
                        "--toastBarBackground": "#2F855A",
                    },
                });
            } else {
                console.log("Authentication failed");
                toast.push("Wrong password!", {
                    theme: {
                        "--toastColor": "mintcream",
                        "--toastBackground": "rgba(255,0,0,0.9)",
                        "--toastBarBackground": "red",
                    },
                });
            }
        } catch (error) {
            console.error("An error occurred during authentication:", error);
            toast.push("Something went wrong... Try again?", {
                theme: {
                    "--toastColor": "mintcream",
                    "--toastBackground": "rgba(255,0,0,0.9)",
                    "--toastBarBackground": "red",
                },
            });
        }
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
