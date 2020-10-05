class SSE {
    constructor(sseUrl, callback) {
        this.sseUrl = sseUrl
        this.callback = callback
    }
    run() {
        let es = new EventSource(this.sseUrl);
        es.onerror = (err) => {console.log(err)}
        es.addEventListener('message', event => {
            let data = JSON.parse(event.data);
            this.callback(data)
        }, false);
    }
}

export default SSE