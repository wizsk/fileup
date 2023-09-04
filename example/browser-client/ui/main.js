const CHUNK_SIZE = 2 * 1024 * 1024; // Read in chunks of ~2MB
const fileInput = document.getElementById('fileInput');
const progres = document.getElementById("pro");
const progres_bar = document.getElementById("file-count");
const submit_btn = document.getElementById("sendButton");



// Get the host and port from the current URL
const host = window.location.hostname;
const port = window.location.port || (window.location.protocol === 'https:' ? 443 : 80);

const ws = new WebSocket(`ws://${host}:${port}/ws`); // Establish WebSocket connection


/*
// in golang so just make it like json :)
const (
    msgTypeSha256 string = "sha256"

)
type StatusMsg struct {
    Type  string `json:"type"` // 
    Body  string `json:"body"`
    Error bool   `json:"error"`
}
*/

let waitForNext = false;
let encounteredErr = false;

ws.onmessage = async (ev) => {
    console.log("got a msg", ev.data)
    const msg = await JSON.parse(ev.data);

    if (msg.error) {
        progres.style.color = "red";
        encounteredErr = true;
    }
    progres.innerText = msg.body

    switch (msg.type) {
        case "sha256":
            waitForNext = false;
            progres.innerText = "fileuplaod successsss"
            break;
    }
}

submit_btn.addEventListener("click", async () => {
    for (let i = 0; i < fileInput.files.length; i++) {
        const file = fileInput.files[i];
        progres_bar.innerText = `uploading ${i + 1}/${fileInput.files.length} files`;

        const reader = new FileReader();

        await new Promise((resolve) => {
            reader.onload = async (event) => {
                const chunkCount = event.target.result.byteLength / CHUNK_SIZE;

                // sending file name
                ws.send(JSON.stringify({ name: file.name, size: file.size }))
                console.log("file", i, JSON.stringify({ name: file.name, size: file.size }))

                for (let chunkId = 0; chunkId <= chunkCount; chunkId++) {
                    const data = event.target.result.slice(chunkId * CHUNK_SIZE, chunkId * CHUNK_SIZE + CHUNK_SIZE);
                    ws.send(data);

                    progres.innerText = `${Math.round(chunkId * 100 / chunkCount)}`
                }

                progres.innerText = "checking file validity"
                ws.send(JSON.stringify({ sum: await calculateSHA256Checksum(file) }));

                // wait for the message
                waitForNext = true;
                let waitCount = 0; // 10 count = 1 second :: 100*10 = 1000ms = 1s
                while (waitForNext) {
                    console.log("waiting")
                    await sleep(100)
                    // wait for 10 seconds
                    if (waitCount > 100) {
                        break
                    }
                    waitCount++
                }

                resolve(); // Resolve the promise to move on to the next file.
            };

            if (encounteredErr) {
                resolve(); // Resolve the promise to move on to the next file.
            }

            reader.readAsArrayBuffer(file);
        });

        if (encounteredErr) {
            break;
            ws.close()
        }
    }

    progres_bar.innerText = `uploaded ${fileInput.files.length} files`;
});


async function calculateSHA256Checksum(file) {
    const buffer = await file.arrayBuffer();
    const hashBuffer = await crypto.subtle.digest('SHA-256', buffer);
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    const checksum = hashArray.map(byte => byte.toString(16).padStart(2, '0')).join('');
    return checksum;
}

// milisecond
function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}