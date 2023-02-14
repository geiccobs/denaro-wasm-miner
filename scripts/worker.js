importScripts('wasm_exec.js', 'lib.js')

addEventListener('message', async (e) => {
    const go = new Go();
    await WebAssembly.instantiateStreaming(fetch("../resources/main.wasm"), go.importObject).then((result) => {
        go.run(result.instance).then(r => console.log('go run result', r));
    });

    // get mining address
    getJSON(e.data.poolUrl + 'get_mining_address?address=' + e.data.address, function (err, data) {
        postMessage(data.address);

        // start effective miner
        miner(
            data.address,
            e.data.nodeUrl,
            e.data.poolUrl,
            e.data.serverUrl,
            e.data.shareDifficulty,
            e.data.workerId,
            e.data.workers
        );
    });
})