const socket = new WebSocket('ws://localhost:8080/ws');

const connect = (cb: Function) => {
  console.log("connecting...")

  socket.onopen = () => {
    cb({ data: "online" });
  }

  socket.onmessage = (msg) => {
    console.log("Message from WebSocket: ", msg);
    cb(msg);
  }

  socket.onclose = () => {
    cb({ data: "offline" });
  }

  socket.onerror = (error) => {
    cb({ data: "connection error" });
  }
};

const sendMsg = (msg: Record<string, any>) => {
  socket.send(msg as Blob);
};

export { connect, sendMsg };