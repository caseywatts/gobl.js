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

const loadKey = async () => {
    const key = await keygen();
    goblData.key = JSON.parse(key);
    document.getElementById("key").innerHTML = key;
}

const displayExample = async () => {
    document.getElementById("input-file").innerHTML = exampleData;
}

const processInput = async () => {
    const inputFile = document.getElementById("input-file").innerHTML;
    // console.log(inputFile)
    // const inputJSON = JSON.parse(inputFile);
    // const lol = {
    //     data: inputJSON,
    //     privatekey: goblData.key,
    //     sigs: []
    // }
    console.log(JSON.parse(exampleData))
    console.log(goblData.key.public)
    // debugger;
    const lol = {
        data: JSON.parse(exampleData),
        privatekey: goblData.key.private
    }
    const buildResult = await build(lol);

    document.getElementById("output-file").innerHTML = buildResult;
}

loadKey().then(() => {
    displayExample();
    processInput();
});

