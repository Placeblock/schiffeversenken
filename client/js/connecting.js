
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

export function startAnimation() {
    changeText()
    animationInterval = setInterval(changeText, 5000);
    dotInterval = setInterval(changeDots, 300);
}

export function stopAnimation() {
    clearInterval(animationInterval);
    clearInterval(dotInterval);
}

let i = 0;
function changeText() {
    console.log("TEST")
    textElement.innerText = texts[i]
    i = (i+1)%texts.length;
}

let dots = 0;
function changeDots() {
    textElement.innerText = texts[i]+".".repeat(dots)
    dots = (dots+1)%4
}