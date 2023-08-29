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
        console.error(msg.body)
        return
    }

    console.log(msg.body)
}


fileInput.addEventListener('change', async () => {
    console.log(fileInput)
    const file = fileInput.files[0];

    const reader = new FileReader();

    reader.onload = async (event) => {
        const chunkCount = event.target.result.byteLength / CHUNK_SIZE;
        console.log(file)
        ws.send(file.name)

        for (let chunkId = 0; chunkId <= chunkCount; chunkId++) {
            const data = event.target.result.slice(chunkId * CHUNK_SIZE, chunkId * CHUNK_SIZE + CHUNK_SIZE);
            progres.innerText = `${chunkId}`;
            ws.send(data);
        }
        ws.send(await calculateSHA256Checksum(file))
    };

    reader.readAsArrayBuffer(file);
});

// async function calculateFileChecksum() {
//     const file = fileInput.files[0];

//     if (file) {
//         const checksum = await calculateSHA256Checksum(file);
//         // document.getElementById("checksumResult").textContent = "Checksum: " + checksum;
//     } else {
//         // document.getElementById("checksumResult").textContent = "Please select a file.";
//     }
// }

async function calculateSHA256Checksum(file) {
    const buffer = await file.arrayBuffer();
    const hashBuffer = await crypto.subtle.digest('SHA-256', buffer);
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    const checksum = hashArray.map(byte => byte.toString(16).padStart(2, '0')).join('');
    return checksum;
}