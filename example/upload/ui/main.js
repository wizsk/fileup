const CHUNK_SIZE = 2 * 1024 * 1024; // Read in chunks of 2 MB
const fileInput = document.getElementById('fileInput');
const progres = document.getElementById("pro")



// Get the host and port from the current URL
const host = window.location.hostname;
const port = window.location.port || (window.location.protocol === 'https:' ? 443 : 80);

const ws = new WebSocket(`ws://${host}:${port}/ws`); // Establish WebSocket connection


ws.onmessage = async (ev) => {
    const msg = await JSON.parse(ev.data);

    if (msg.is_error) {
        progres.style.color = "red";
    }

    console.log(msg.body)
    progres.innerText = msg.body
}


fileInput.addEventListener('change', async () => {
    // for (let i = 0; i< fileInput.files.lenth; )
    for (let i = 0; i < fileInput.files.length; i++) {
        const file = fileInput.files[i];
        console.log("file", file.name)

        const reader = new FileReader();

        reader.onload = async (event) => {
            const chunkCount = event.target.result.byteLength / CHUNK_SIZE;

            // sending file name
            ws.send(JSON.stringify({ file_name: file.name }))

            const b = new ArrayBuffer(8)
            ws.send(b)

            // for (let chunkId = 0; chunkId <= chunkCount; chunkId++) {
            //     const data = event.target.result.slice(chunkId * CHUNK_SIZE, chunkId * CHUNK_SIZE + CHUNK_SIZE);
            //     progres.innerText = `${chunkId}`;
            //     ws.send(data);
            //     progres.innerText = `${Math.round(chunkId * 100 / chunkCount)}`
            //     await sleep(10)
            // }

            // progres.innerText = "ceheking file validity"
            // ws.send(JSON.stringify({ checksum: await calculateSHA256Checksum(file) }));
            ws.send(JSON.stringify({ checksum: "foo" }));
        };

        reader.readAsArrayBuffer(file);
    }
});

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}


async function calculateSHA256Checksum(file) {
    const buffer = await file.arrayBuffer();
    const hashBuffer = await crypto.subtle.digest('SHA-256', buffer);
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    const checksum = hashArray.map(byte => byte.toString(16).padStart(2, '0')).join('');
    return checksum;
}