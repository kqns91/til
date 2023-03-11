import React, {useEffect, useRef} from "react";
import { useParams } from "react-router-dom";

const Room = () => {
  const { roomID } = useParams();

  console.log(roomID)

  useEffect(() => {
    const ws = new WebSocket(
      `ws://localhost:8000/join?roomID=${roomID}`
    );

    ws.addEventListener("open", () => {
      ws.send(JSON.stringify({join: "true"}));
    })
  })

  return (
    <div>
       <video autoPlay controls={true}></video>
       <video autoPlay controls={true}></video>
    </div>
  )
}

export default Room;