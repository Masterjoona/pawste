<script>
	import Header from "../lib/ui/header.svelte"
	import Footer from '../lib/ui/footer.svelte';
	import './styles.css';
	import { onMount } from "svelte";
	let display = false;
	onMount(() => {
		setTimeout(() => {
			display = true
		}, 50);
	})
</script>


<div id ="container" class:hidden={!display}>
	<div id ="header"><Header/></div>
	
	<main>
		<slot />
	</main>
	
	<div id ="footer"><Footer/></div>
</div>

<style>
	#container {
		display: flex;
		flex-direction: column;
		align-items: center;
	}
	#footer {
		position: absolute;
		bottom: 0px;
		right: 0px;
	}
	main {
		width: 80vw;
        height: 80vh;
	}
	.hidden {
		opacity: 0;
	}
</style>