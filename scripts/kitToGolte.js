import { promises as fs } from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const srcDir = path.join(__dirname, 'svelte-dev', 'src', 'routes');
const destDir = path.join(__dirname, 'web', 'page');

(async () => {
  try {
    await fs.mkdir(destDir, { recursive: true });

    const files = await fs.readdir(srcDir, { withFileTypes: true });

    const movePromises = files.map(async (file) => {
      if (file.isDirectory()) {
        const folderName = file.name;
        const srcFilePath = path.join(srcDir, folderName, '+page.svelte');
        const destFilePath = path.join(destDir, `${folderName}.svelte`);
        try {
          // For some husk reason fs.copyFile errors and as everyone knows switching to secondary is faster than reloading, so here's my sidearm copy
          const data = await fs.readFile(srcFilePath);
          await fs.writeFile(destFilePath, data);
        } catch (err) {
          console.error(`Error moving file: ${srcFilePath} to ${destFilePath}`, err);
        }
      }
    });

    await Promise.all(movePromises);

    console.log('done :3');
  } catch (err) {
    console.error('some husk error', err);
  }
})();