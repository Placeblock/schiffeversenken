class MessageTarget extends EventTarget {
    constructor() {
        super()
    }
}

const messageTarget = new MessageTarget()
export {messageTarget}

const socket = new WebSocket("ws://localhost:4195");

socket.onmessage = (event) => {
    const message = JSON.parse(event.data);
    const e = new CustomEvent(message.action, {detail: message.data})
    messageTarget.dispatchEvent(e)
    console.log(message);
};

export function sendMessage(action, data) {
    socket.send(JSON.stringify({action, data}))
}