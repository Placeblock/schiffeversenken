export const State = {
    Connecting: "connecting",
    Pool: "pool",
    Building: "building",
    Playing: "playing",
    Ended: "ended"
}

class StateTarget extends EventTarget {
    constructor() {
        super()
    }
}

const stateTarget = new StateTarget()
export { stateTarget }

let state = State.Connecting;
export function setState(newState) {
    if (newState == state) return
    const section = document.getElementById(state)
    section.style.display = "none"
    state = newState
    const newSection = document.getElementById(state)
    newSection.style.display = "flex"
    
    const stateEvent = new CustomEvent(newState)
    const event = new CustomEvent("state", {detail: {state: newState}})
    stateTarget.dispatchEvent(stateEvent)
    stateTarget.dispatchEvent(event)
}