import { setRoomID } from "./room.js";
import { State, setState } from "./state.js";
import { messageTarget } from "./ws.js";

const textElement = document.getElementById("loading-text")

const texts = [
    "Spitze Bleistifte",
    "Zeichne das Spielfeld",
    "Baue den Sichtschutz auf",
    "Bereite das Spiel vor",
    "Fälle den Baum",
    "Stelle Papier her",
    "Drucke Kästchen auf das Papier",
]

let animationInterval;
let dotInterval;

function startAnimation() {
    changeText()
    animationInterval = setInterval(changeText, 5000);
    dotInterval = setInterval(changeDots, 300);
}

function stopAnimation() {
    clearInterval(animationInterval);
    clearInterval(dotInterval);
}

let i = 0;
function changeText() {
    textElement.innerText = texts[i]
    i = (i+1)%texts.length;
}

let dots = 0;
function changeDots() {
    textElement.innerText = texts[i]+".".repeat(dots)
    dots = (dots+1)%4
}

messageTarget.addEventListener("ROOM", (data) => {
    setRoomID(data.detail)
    stopAnimation()
    setState(State.Pool)
})

startAnimation()