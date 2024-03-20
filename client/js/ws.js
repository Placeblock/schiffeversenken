class MessageTarget extends EventTarget {
    constructor() {
        super()
    }
}

const messageTarget = new MessageTarget()
export {messageTarget}

const socket = new WebSocket("wss://battleship.codelix.de/wss");

socket.onopen = () => {
    const urlParams = new URLSearchParams(window.location.search);
    const id = urlParams.get("room")
    if (id != null && urlParams.get("shared") != null) {
        sendMessage("JOIN", id)
    }
}

socket.onmessage = (event) => {
    const message = JSON.parse(event.data);
    const e = new CustomEvent(message.action, {detail: message.data})
    messageTarget.dispatchEvent(e)
    console.log(message);
};

export function sendMessage(action, data) {
    socket.send(JSON.stringify({action, data}))
}