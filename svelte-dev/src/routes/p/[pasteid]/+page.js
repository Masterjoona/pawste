import { error } from '@sveltejs/kit';

function getPaste(pasteId) {
	// pretend like this is fetching data
	return {
		id: pasteId,
		title: 'Epic paste',
		textContent: 'owo meow meow nya purr',
		readCount: 69,
		bytes: 1337
	};
}
/** @type {import('./$types').PageLoad} */
export function load({ params }) {
	return getPaste(params.pasteid);

	error(404, 'Not found');
}
