export const State = {
    Connecting: "connecting",
    Pool: "pool",
    Building: "building",
    Playing: "playing"
}

let state = State.Connecting;
export function setState(newState) {
    const section = document.getElementById(state)
    section.style.display = "none"
    state = newState
    const newSection = document.getElementById(state)
    newSection.style.display = "flex"
}