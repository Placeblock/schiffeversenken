import { State, setState } from "./state";

class MessageTarget extends EventTarget {
    constructor() {
        super()
    }
}

const messageTarget = new MessageTarget()
export { messageTarget }

let socket
function connect() {
    socket = new WebSocket("wss://battleship.codelix.de/wss");

    socket.onopen = () => {
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

    socket.onclose = function (e) {
        setState(State.Connecting);
        console.log('Socket is closed. Reconnect will be attempted in 1 second.', e.reason);
        setTimeout(function () {
            connect();
        }, 1000);
    };

}

connect();

export function sendMessage(action, data) {
    socket.send(JSON.stringify({ action, data }))
}