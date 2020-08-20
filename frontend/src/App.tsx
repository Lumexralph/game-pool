import React, { useEffect, useState } from "react";
import { Button, Table, Form, InputNumber, Input, Alert, Badge, Tag } from "antd";
import styled from "@emotion/styled";

import "./App.css";
import { connect, sendMsg } from "./api";

const columns = [
  {
    title: "Position",
    dataIndex: "key",
    key: "key",
  },
  {
    title: "name",
    dataIndex: "name",
    key: "name",
  },
  {
    title: "Score",
    dataIndex: "score",
    key: "score",
  },
];

type ScoreBoardProps = {
  dataSource: {
    key: string;
    name: string;
    score: number;
  }[];
};

const ScoreBoard = ({ dataSource }: ScoreBoardProps) => (
  <Table dataSource={dataSource} columns={columns} />
);

interface PlayerInputProps {
  onFinish: (value: any) => void;
}
const PlayerInput = ({ onFinish }: PlayerInputProps) => (
  <div>
    <Form onFinish={onFinish}>
      <Form.Item name="player">
        <Input placeholder="player name" />
      </Form.Item>
      <Form.Item name="input1" initialValue={1}>
        <InputNumber min={1} max={10} />
      </Form.Item>

      <Form.Item name="input2" initialValue={1}>
        <InputNumber min={1} max={10} />
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit">Play</Button>
      </Form.Item>
    </Form>
  </div>
);

interface CustomAlertProps {
  visible: boolean;
  info: string;
  handleClose: () => void;
}
const CustomAlert = ({ visible, handleClose, info }: CustomAlertProps) => (
  <div>
    {visible ? (
      <Alert message={info} type="success" closable afterClose={handleClose} />
    ) : null}
  </div>
);

type mode = "observe" | "play" | null;

function App() {
  const [history, setHistory] = useState<Record<string, any>[]>([]);
  const [scoreBoard, setScoreBoard] = useState<
    {
      key: string;
      name: string;
      score: number;
    }[]
  >([]);
  const [clientID, setClientID] = useState("");
  const [playerMode, setPlayerMode] = useState<mode>(null);
  const [info, setInfo] = useState("");

  const handlePlayerMode = (mode: mode) => setPlayerMode(mode);

  function onChange(value: string | number | undefined) {
    sendMsg((value as unknown) as Blob);
    console.log("changed", value);
  }

  const onFinish = (values: any) => {
    let msg = { clientID, playerMode: "roundPlay", ...values };

    console.log("Success:", msg);
    msg = JSON.stringify(msg);
    setTimeout(() => sendMsg((msg as unknown) as Record<string, any>), 900);
  };

  const [visible, setVisible] = useState(false);
  const [onlineStatus, setOnlineStatus] = useState(false);

  const handleUpdate = (info: string) => {
    setInfo(info);
    setVisible(true);
  };

  useEffect(() => {
    let deserializedMsg: any;

    connect((msg: Record<string, any>) => {
      switch (true) {
        case msg.data === "offline":
          console.log("status: offline");
          setOnlineStatus(false);
          break;
        case msg.data === "online":
          console.log("status: online");
          setOnlineStatus(true);
          break;
        case msg.data === "connection error":
          console.log("status: cannot establish connection.");
          setOnlineStatus(false);
          break;

        default:
          deserializedMsg = JSON.parse((msg.data as unknown) as string);
          // alert the new info
          handleUpdate(deserializedMsg.info);
          // store the clientID which will be used to play the game
          if (clientID === "") setClientID(deserializedMsg.clientID);

          // store the scoreboard
          if (deserializedMsg.type === "scoreboard") {
            const dataSource = deserializedMsg.scoreboard.map(
              (score: Record<string, any>, index: number) => ({
                key: index + 1,
                name: score.name,
                score: score.totalScore,
              })
            );

            setScoreBoard(dataSource);

            setHistory((prevState) => {
              const newHistory = [...prevState, deserializedMsg];

              console.log(newHistory);
              return newHistory;
            });
          }
      }
    });
  }, [clientID]);

  useEffect(() => {
    setTimeout(() => setVisible(false), 5000);
  }, [visible]);

  const handleClose = () => {
    setVisible(false);
  };

  return (
    <App.Wrapper className="App">
      <section>
        <CustomAlert info={info} visible={visible} handleClose={handleClose} />
      </section>
      <section>
        <p>Play or Watch the game</p>
        <p>
          <Badge status={onlineStatus ? "success" : "default"} text={onlineStatus ? "online" : "offline"} />
          <br />
          {onlineStatus && <Tag color="green">Game In Session</Tag>}
        </p>
        <Button
          type="primary"
          onClick={() => {
            // send signal to play
            sendMsg(
              (JSON.stringify({
                playerMode: "play",
                clientID,
              }) as unknown) as Blob
            );
            handlePlayerMode("play");
          }}
        >
          Play
        </Button>
        <Button onClick={() => handlePlayerMode("observe")}>Watch</Button>
        {playerMode && <ScoreBoard dataSource={scoreBoard} />}
      </section>
      <section>
        {playerMode === "play" && <PlayerInput onFinish={onFinish} />}
      </section>
    </App.Wrapper>
  );
}

App.Wrapper = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  column-gap: 1em;
`;

export default App;
