const getJSON = function (url, callback) {
    const xhr = new XMLHttpRequest();

    xhr.open('GET', url, true);

    xhr.responseType = 'json';
    xhr.onload = function () {
        const status = xhr.status;
        if (status === 200) {
            callback(null, xhr.response);
        } else {
            callback(status, xhr.response);
        }
    };
    xhr.send();
};

/*
    * All the following stuff is used by WASM since we can't do requests directly from it.
    * We need to use JS to do that.
    *
    * Therefore, unused functions reports are expected.
 */
var response = new Map();
function expGetResponse(workerId) {
    let value = response.get(workerId);
    response.delete(workerId);
    return value;
}

function expPostJSON(url, data, workerId) {
    const xhr = new XMLHttpRequest();

    xhr.open('POST', url, false);
    xhr.setRequestHeader('Content-Type', 'application/json');

    xhr.send(JSON.stringify(data));

    if (xhr.status === 200) {
        response.set(workerId, JSON.parse(xhr.response));
    }
}

function expGetJSON(url, workerId) {
    const xhr = new XMLHttpRequest();

    xhr.open('GET', url, false);

    xhr.responseType = 'json';
    xhr.send();

    if (xhr.status === 200) {
        response.set(workerId, xhr.response);
    }
}