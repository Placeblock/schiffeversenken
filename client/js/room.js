import { messageTarget, sendMessage } from "./ws.js";
import { State, setState } from "./state.js"

const roomCodeElement = document.getElementById("room-code")

const roomCodeForm = document.getElementById("room-number-form")

let roomID
export function setRoomID(id) {
    roomID = id
    roomCodeElement.innerText = id
}

messageTarget.addEventListener("ROOM", (data) => {
    setRoomID(data.detail)
    setState(State.Pool)
})

roomCodeForm.onsubmit = ev => {
    ev.preventDefault();
    const formData = new FormData(roomCodeForm);
    const id = formData.get("id")
    sendMessage("JOIN", id)
}

const shareBtn = document.getElementById("share-btn")
shareBtn.onclick = () => {
    const url = new URL(window.location);
    url.searchParams.set("shared", true)
    url.searchParams.set("room", roomID)
    navigator.clipboard.writeText(url.toString());
    alert("In die Zwischenablage kopiert!");
}