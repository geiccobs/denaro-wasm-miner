const workers = [];

function startWorkers() {
    let server = document.getElementById('server').value;

    // create web workers
    for (let i = 0; i < document.getElementById('workers').value; i++) {
        workers.push(new Worker('scripts/worker.js'));
    }

    // we will receive only one message and data will be mining address
    workers[0].addEventListener('message', function (e) {
        setInterval(function () {
            getJSON(server+'getData?address='+e.data, function (err, data) {
                if (err !== null) {
                    alert('Something went wrong: ' + err);
                } else {
                    document.getElementById('hashrate').innerText = 'Hashrate: ' + data.hashrate + 'k hash/s';
                    document.getElementById('shares').innerText = 'Shares: ' + data.shares;
                    document.getElementById('mined_blocks').innerText = 'Mined blocks: ' + data.mined_blocks;
                }
            });
        }, 5_000);
    });

    for (let i = 0; i < workers.length; i++) {
        // start web worker (miner)
        workers[i].postMessage({
            address: document.getElementById('address').value,
            nodeUrl: document.getElementById('node').value,
            poolUrl: document.getElementById('pool').value,
            serverUrl: server,
            shareDifficulty: document.getElementById('share_difficulty').value,
            workerId: i+1,
            workers: workers.length
        });
    }
}

function stopWorkers() {
    for (let i = 0; i < workers.length; i++) {
        workers[i].terminate();
    }
    workers.length = 0;
}

function switchStatus() {
    // set start button as not clickable
    let startElement = document.getElementById('start');

    if (startElement.innerText === 'Start') {
        startWorkers();

        startElement.innerText = 'Stop';
        startElement.style.background = '#800020';
    } else {
        stopWorkers();

        startElement.innerText = 'Start';
        startElement.style.background = '#6a64f1';
    }
}
