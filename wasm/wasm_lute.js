'use strict';

const WASM_URL = '/lute.wasm';
let wasm;

async function init() {
    const go = new Go(); // Defined in wasm_exec.js

    await WebAssembly.instantiateStreaming(fetch(WASM_URL), go.importObject).then((obj)=>{
        wasm = obj.instance;
        go.run(wasm);
    })

    console.log(go.importObject);
    console.log('wasm loaded');
}

init().then(()=>{
	console.log(Md2Html('**Lute** - A structured markdown engine.'));
})

