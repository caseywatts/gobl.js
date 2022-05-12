import { keygen, build } from "./gobl.js";

// assigning these to the global namespace for cypress tests
window.gobl = {};
window.gobl.keygen = keygen;
window.gobl.build = build;

let goblData = {};

const exampleInputs = {};
exampleInputs.simpleMessage = `{
    "doc": { 
        "$schema": "https://gobl.org/draft-0/note/message",
        "title": "Test Message",
        "content": "test content"
    }
}`;
exampleInputs.noSchema = `{
    "doc": { 
        "title": "Test Message",
        "content": "test content"
    }
}`;

const exampleData = exampleInputs.simpleMessage;

const generateAndDisplayKey = async () => {
    const key = await keygen();
    goblData.key = JSON.parse(key);
    document.getElementById("key").value = key;
}

const displayExampleInputFile = async () => {
    document.getElementById("input-file").value = exampleData;
}

const processInputFile = async () => {
    const inputFile = document.getElementById("input-file").value;

    const buildData = {
        data: JSON.parse(inputFile),
        privatekey: goblData.key.private
    }

    try { 
        const buildResult = await build(buildData);
        document.getElementById("output-file").value = buildResult;
        updateStatus("success");
    } catch (e) {
        document.getElementById("output-file").value = "";
        updateStatus("error", e);
    }
}


const markSuccess = (el) => {
    el.classList.remove("bg-red-200")
    el.classList.add("bg-green-200")
}

const markError = (el) => {
    el.classList.add("bg-red-200")
    el.classList.remove("bg-green-200")
}

const updateStatus = async (type, message) => {
    const statusEl = document.getElementById("status");
    if (type === "success") {
        statusEl.innerHTML = "Success!"
        markSuccess(statusEl)
    } else { // error case
        statusEl.innerHTML = `Error: ${message}`
        markError(statusEl)
    }
}

await generateAndDisplayKey();
// await displayExampleInputFile();
await processInputFile();

document.getElementById("input-file").oninput = function updateOnInputFileChange () {
    processInputFile();
};