import { keygen, build } from "./gobl.js";

window.gobl = {};
window.gobl.keygen = keygen;
window.gobl.build = build;

// const result = await keygen();
// console.log(`RESULT: ${result}`);

// try {
//     const buildResult = await build({
//         "data": {},
//         "privateykey": {},
//     });
//     console.log(`BUILD RESULT: ${build_result}`);
// } catch (e) {
//     console.log("BUILD ERROR: " + e)
// };

// const result2 = await keygen();
// console.log(`RESULT2: ${result2}`);

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
    const buildResult = await build(buildData);

    document.getElementById("output-file").value = buildResult;
}

(async function loadExample() {
    await generateAndDisplayKey();
    await displayExampleInputFile();
    await processInputFile();
})()

document.getElementById("input-file").oninput = function updateOnInputFileChange () {
    processInputFile();
};