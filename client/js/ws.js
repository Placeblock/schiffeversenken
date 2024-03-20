import { State, setState } from "./state.js";

class MessageTarget extends EventTarget {
    constructor() {
        super()
    }
}

const messageTarget = new MessageTarget()
export { messageTarget }

let socket
function connect() {
    console.log("Connecting...")
    socket = new WebSocket("wss://battleship.codelix.de/wss");

    socket.onopen = () => {
        console.log("Opened")
        const urlParams = new URLSearchParams(window.location.search);
        const id = urlParams.get("room")
        if (id != null && urlParams.get("shared") != null) {
            sendMessage("JOIN", id)
        }
    }

    socket.onmessage = (event) => {
        const message = JSON.parse(event.data);
        const e = new CustomEvent(message.action, { detail: message.data })
        messageTarget.dispatchEvent(e)
        console.log(message);
    };

    socket.onclose = (e) => {
        setState(State.Connecting);
        console.log('Socket is closed. Reconnect will be attempted in 1 second.', e.reason);
        setTimeout(connect, 1000);
    };

}

connect();

export function sendMessage(action, data) {
    socket.send(JSON.stringify({ action, data }))
}