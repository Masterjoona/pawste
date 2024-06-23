<script>
	import Header from '$lib/ui/Header.svelte';
	import Footer from '$lib/ui/footer.svelte';
	import './styles.css';
</script>

<div id ="container">
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
</style>