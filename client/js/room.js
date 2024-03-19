import { sendMessage } from "./ws.js";

const roomCodeElement = document.getElementById("room-code")

const roomCodeForm = document.getElementById("room-number-form")

export function setRoomID(id) {
    roomCodeElement.innerText = id
}

roomCodeForm.onsubmit = ev => {
    ev.preventDefault();
    const formData = new FormData(roomCodeForm);
    const id = formData.get("id")
    sendMessage("JOIN", id)
}