import * as animation from "./connecting.js";
import { setRoomID } from "./room.js";
import { State, setState } from "./state.js";

animation.startAnimation()
const socket = new WebSocket("ws://localhost:4195",);

socket.onmessage = (event) => {
    const message = JSON.parse(event.data);
    switch (message.action) {
        case "ROOM":  
            animation.stopAnimation()
            setState(State.Pool)
            setRoomID(message.data)
            break;
        default:
            break;
    }
    console.log(message);
};