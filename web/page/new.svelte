<script>
    import { toast } from "@zerodevx/svelte-toast";
    import { truncateFilename } from "../lib/utils.js";
    import "../styles/paste.css";

    let selectedExpiration = "1w";
    let selectedBurnAfter = "0";
    let selectedSyntax = "none";
    let selectedPrivacy = "public";
    let content = "";
    let attachedFiles = [];
    let imageSources = [];
    let password = "";

    function handleAttachFiles(event) {
        const files = event.target.files;
        for (let file of files) {
            attachedFiles = [...attachedFiles, file];
            if (file.type.startsWith("image/")) {
                const reader = new FileReader();
                reader.onload = (e) => {
                    imageSources = [...imageSources, e.target.result];
                };
                reader.readAsDataURL(file);
            } else {
                imageSources = [...imageSources, null];
            }
        }
    }

    async function handleSave() {
        if (content === "" && attachedFiles.length === 0) {
            toast.push("You must provide content or attach files!", {
                theme: {
                    "--toastColor": "mintcream",
                    "--toastBackground": "rgba(255,0,0,0.9)",
                    "--toastBarBackground": "red",
                },
            });
            return;
        }

        const encrypted =
            selectedPrivacy === "private" || selectedPrivacy === "secret";

        if (encrypted && password === "") {
            toast.push("You must provide a password for encrypted pastes!", {
                theme: {
                    "--toastColor": "mintcream",
                    "--toastBackground": "rgba(255,0,0,0.9)",
                    "--toastBarBackground": "red",
                },
            });

            return;
        }
        /*const data = {
            expiration: selectedExpiration,
            burnAfter: selectedBurnAfter,
            syntax: selectedSyntax,
            privacy: selectedPrivacy,
            content: content,
            files: attachedFiles,
            password: encrypted ? password : null,
        };
        console.log("Data saved:", data);
        alert("Data saved successfully!");*/
        const formData = new FormData();
        formData.append("expire", selectedExpiration);
        formData.append("burnafter", selectedBurnAfter);
        formData.append("syntax", selectedSyntax);
        formData.append("privacy", selectedPrivacy);
        formData.append("content", content);
        formData.append("password", password);
        console.log(attachedFiles);
        for (let file of attachedFiles) {
            formData.append("files[]", file);
        }

        const response = await (
            await fetch("/p", {
                method: "POST",
                body: formData,
            })
        ).json();

        if (response?.error) {
            toast.push(response.error, {
                theme: {
                    "--toastColor": "mintcream",
                    "--toastBackground": "rgba(255,0,0,0.9)",
                    "--toastBarBackground": "red",
                },
            });
        } else {
            toast.push("Paste saved successfully! Redirecting...", {
                theme: {
                    "--toastColor": "mintcream",
                    "--toastBackground": "rgba(0,255,0,0.9)",
                    "--toastBarBackground": "green",
                },
            });

            setTimeout(() => {
                window.location.href =
                    `/p/${response.pasteName}` + (encrypted ? "/auth" : "");
            }, 500);
        }
    }

    function removeFile(index) {
        attachedFiles = attachedFiles.filter((_, i) => i !== index);
    }
    function handlePrivacyChange(event) {
        selectedPrivacy = event.target.value;
        const passwordField = document.getElementById("password-field");
        if (selectedPrivacy === "private" || selectedPrivacy === "secret") {
            passwordField.style.display = "block";
        } else {
            passwordField.style.display = "none";
        }
    }
</script>

<div id="container">
    <div class="card">
        <div class="options">
            <div>
                <label for="expiration">Expiration:</label>
                <select id="expiration" bind:value={selectedExpiration}>
                    <option value="never">Never</option>
                    <option value="1h">1 Hour</option>
                    <option value="6h">6 Hours</option>
                    <option value="1d">1 Day</option>
                    <option value="3d">3 Days</option>
                    <option value="1w">1 week</option>
                </select>
            </div>
            <div>
                <label for="burn-after">Burn After:</label>
                <select id="burn-after" bind:value={selectedBurnAfter}>
                    <option value="0">Never</option>
                    <option value="1">1 View</option>
                    <option value="10">10 Views</option>
                    <option value="100">100 Views</option>
                    <option value="1000">1000 Views</option>
                </select>
            </div>
            <div>
                <label for="syntax">Syntax:</label>
                <select id="syntax" bind:value={selectedSyntax}>
                    <option value="none">None</option>
                </select>
            </div>
            <div>
                <label for="privacy">Privacy:</label>
                <select
                    id="privacy"
                    bind:value={selectedPrivacy}
                    on:change={handlePrivacyChange}>
                    <option value="public">Public</option>
                    <option value="unlisted">Unlisted</option>
                    <option value="readonly">Read-only</option>
                    <option value="private">Private</option>
                    <option value="secret">Secret</option>
                </select>
            </div>
            <div id="password-field">
                <label for="password">Password:</label>
                <input
                    type="password"
                    id="password"
                    on:input={(e) => (password = e.target.value)} />
            </div>
        </div>

        <textarea placeholder="Pawste away" bind:value={content}></textarea>
        <div class="buttons">
            <input
                type="file"
                multiple
                on:change={handleAttachFiles}
                style="display: none;"
                id="file-input" />
            <button
                on:click={() => document.getElementById("file-input").click()}
                >Attach Files</button>
            <button on:click={handleSave}>Save</button>
        </div>
        <div class="file-list">
            {#each attachedFiles as file, index}
                <div class="file-item">
                    {#if file.type.startsWith("image/")}
                        <img
                            src={imageSources[index]}
                            alt={file.name}
                            class="thumbnail" />
                    {/if}
                    <span>
                        {truncateFilename(file.name)} - {(
                            file.size / 1024
                        ).toFixed(2)} KB
                    </span>
                    <button on:click={() => removeFile(index)}>Remove</button>
                </div>
            {/each}
        </div>
    </div>
</div>

<style>
    .options {
        display: flex;
        justify-content: space-evenly;
        flex-direction: column;
        gap: 10px;
        margin-bottom: 10px;
    }

    .options div {
        display: flex;
        align-items: center;
    }

    label {
        font-size: var(--font-size);
        color: white;
    }

    select {
        background-color: #1e1e1e;
        color: white;
        border: 1px solid #444;
        padding: 5px;
        border-radius: 5px;
        font-family: var(--main-font);
        font-size: var(--font-size);
    }

    #password-field {
        margin-top: -0.1%;
        display: none;
    }

    #password {
        width: 40%;
        padding: 5px;
        border-radius: 5px;
        border: 1px solid #444;
        background-color: #1e1e1e;
        color: white;
        font-family: var(--main-font);
        font-size: var(--font-size);
    }

    @media (min-width: 600px) {
        .options {
            flex-direction: row;
            flex-wrap: wrap;
            gap: 20px;
        }

        textarea {
            height: 400px;
        }
    }
</style>
