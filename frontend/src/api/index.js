const loc = window.location;
let wsProtocol = loc.protocol === "https:" ? "wss" : "ws";
const socketUrl = `${wsProtocol}://${loc.host}/ws`;

let socket = new WebSocket(socketUrl);

const connect = (cb) => {
  console.log("Connecting to WebSocket:", socketUrl);

  socket.onopen = () => {
    console.log("Successfully Connected");
  };

  socket.onmessage = (msg) => {
    console.log("Message from WebSocket:", msg);
    cb(msg);
  };

  socket.onclose = (event) => {
    console.log("Socket Closed Connection:", event);
  };

  socket.onerror = (error) => {
    console.error("Socket Error:", error);
  };
};

const sendMsg = (msg) => {
  if (socket.readyState === WebSocket.OPEN) {
    console.log("Sending message:", msg);
    socket.send(msg);
  } else {
    console.warn("WebSocket not open. Message not sent:", msg);
  }
};

export { connect, sendMsg };
