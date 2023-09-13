const CHUNK_SiZE = 1024 * 1024 * 5; // 5MB

const UPLOAD_URL = "/api/upload"

const fileInput = document.getElementById("fileInput");
const fileSubmit = document.getElementById("fileSubmit");
const fileProgress = document.getElementById("progress");
const fileProgressSend = document.getElementById("progress-send");

let isUploading = false;

fileSubmit.addEventListener("click", async () => {
    if (isUploading) {
        console.log("already uploaing")
        return
    }

    if (fileInput.files.length < 1) {
        console.log("select a file");
        return
    }

    await uploadFile();
})


async function uploadFile() {
    isUploading = true;
    fileProgressSend.innerText = `send 0/${fileInput.files.length}`
    for (let i = 0; i < fileInput.files.length; i++) {
        let file = fileInput.files[i];
        let fileURL = `${UPLOAD_URL}/${file.name}`
        const uuid = generateUUID();
        await upload(file, fileURL, uuid);
        fileProgressSend.innerText = `send ${i + 1}/${fileInput.files.length}`;
        console.log("done", file.name);
    }
    isUploading = false;
}

async function upload(file, fileURL, uuid) {
    try {
        await fetch(fileURL, {
            method: 'POST',
            headers: {
                "Upload-Offset": "",
                "Upload-Size": file.size,
                "Upload-UUID": uuid,
            },

        })
    } catch (err) {
        console.error("while uploading file", file.name, err)
        fileProgress.innerText = err;
        return
    }

    const chuckCount = Math.ceil(file.size / CHUNK_SiZE);
    for (let chuckId = 0; chuckId < chuckCount; chuckId++) {
        // calculation mistakes
        const offset = chuckId * CHUNK_SiZE;
        const readUntil = (chuckId * CHUNK_SiZE) + CHUNK_SiZE;
        const data = file.slice(offset, readUntil);

        // it's js so unpridictable :)
        // how must time wated doint this crap :)
        if (data.size === 0) break;

        try {
            await fetch(fileURL, {
                method: 'PUT',
                headers: {
                    "Upload-Offset": offset,
                    "Upload-Size": file.size,
                    "Upload-UUID": uuid,
                    "Content-Type": "application/offset+octet-stream",
                },
                body: data,
            })
        } catch (err) {
            console.error("while uploading file", file.name, err)
            fileProgress.innerText = err;
            return
        }

        const msg = `send ${null}/${fileInput.files.length} <br>curret: ${Math.round((chuckId / chuckCount) * 100)}% ${file.name}`;
        fileProgress.innerHTML = msg;
        console.log(msg.replace("<br>", "\n"));

    }

    // last request
    // rename the file 
    try {
        await fetch(fileURL, {
            method: 'PATCH',
            headers: {
                "Upload-Size": file.size,
                "Upload-UUID": uuid,
                "Sha256": "", // TODO: add this
            },
        })
    } catch (err) {
        console.error("while uploading file", file.name, err)
        fileProgress.innerText = err;
        return
    }
}

// Generate a random UUID
function generateUUID() {
    const cryptoObj = window.crypto || window.msCrypto; // For cross-browser compatibility
    if (cryptoObj && cryptoObj.getRandomValues) {
        // Use a cryptographically strong random number generator if available
        const buffer = new Uint16Array(8);
        cryptoObj.getRandomValues(buffer);
        return `${toHex(buffer[0], 4)}${toHex(buffer[1], 4)}${toHex(buffer[2], 4)}${toHex(buffer[3], 4)}${toHex(buffer[4], 4)}${toHex(buffer[5], 4)}${toHex(buffer[6], 4)}${toHex(buffer[7], 4)}`;
        // }
    } else {
        // Fallback to a less secure method
        return 'xxxxxxxxxxxx4xxxyxxxxxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
            const r = Math.random() * 16 | 0;
            const v = c === 'x' ? r : (r & 0x3 | 0x8);
            return v.toString(16);
        });
    }
}

function toHex(value, width) {
    const hex = value.toString(16);
    return '0'.repeat(width - hex.length) + hex;
}

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}
