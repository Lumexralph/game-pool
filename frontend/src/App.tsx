import React, { useEffect, useState } from "react";
import { Button, Table, Form, InputNumber, Input, Alert, Badge, Tag, Modal } from "antd";
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
      <Form.Item name="input1" initialValue={0}>
        <InputNumber min={1} max={10} />
      </Form.Item>

      <Form.Item name="input2" initialValue={0}>
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

  const onFinish = (values: any) => {
    let msg = { clientID, playerMode: "roundPlay", ...values };
    msg = JSON.stringify(msg);
    setTimeout(() => sendMsg((msg as unknown) as Record<string, any>), 900);
  };

  const [visible, setVisible] = useState(false);
  const [onlineStatus, setOnlineStatus] = useState(false);
  const [gameInSession, setGameInSession] = useState(false);
  const [visibleRoundInfo, setVisibleRoundInfo] = useState(false);
  const [roundInfo, setRoundInfo] = useState("round info");

  const handleUpdate = (info: string) => {
    setInfo(info);
    setVisible(true);
  };

  const displayWinner = (player: Record<string, any>) => {
    const modal = Modal.success({
      title: 'Game Winner',
      content: `The winner is ${player.name} | Total Score : ${player.totalScore}.`,
    });
    setTimeout(() => {
      modal.destroy();
    }, 5 * 1000);
  }

  useEffect(() => {
    let deserializedMsg: any;

    connect((msg: Record<string, any>) => {
      switch (true) {
        case msg.data === "offline":
          setOnlineStatus(false);
          break;
        case msg.data === "online":
          setOnlineStatus(true);
          break;
        case msg.data === "connection error":
          setOnlineStatus(false);
          break;

        default:
          deserializedMsg = JSON.parse((msg.data as unknown) as string);
          // alert the new info
          if (deserializedMsg.type === "game-info") {
            handleUpdate(deserializedMsg.info);
          }

          // TODO: Wrap al the if in a switch statement
          // store the clientID which will be used to play the game
          if (clientID === "") setClientID(deserializedMsg.clientID);

          // when the game starts
          if (deserializedMsg.type === "game-start") {
            setGameInSession(true);
          }

          if (deserializedMsg.type === "game-end") {
            setGameInSession(false);
          }

          if (deserializedMsg.type === "game-winner") {
            displayWinner(deserializedMsg.Winner);
          }

          if (deserializedMsg.type === "player-wait") {
            handleUpdate(deserializedMsg.info);
          }
          // store the scoreboard
          if (deserializedMsg.type === "scoreboard") {
            const dataSource = deserializedMsg.scoreboard.map(
              (score: Record<string, any>, index: number) => ({
                key: index + 1,
                name: score.name,
                score: score.totalScore,
              })
            );

            if (!visibleRoundInfo) setVisibleRoundInfo(true);
            setRoundInfo(deserializedMsg.info);
            setScoreBoard(dataSource);
          }
      }
    });
  }, [clientID, visibleRoundInfo]);

  useEffect(() => {
    setTimeout(() => setVisible(false), 5000);
  }, [visible]);

  const handleClose = () => {
    setVisible(false);
  };

  const gameInstruction = () => {
    Modal.info({
      title: "Game Instructions",
      content: (
        <div>
          <p>You can choose to play or just watch the game. To play the game, click the play button and
          click watch to observe the game.
          </p>
          <p>Should you choose to play, this is how the game works: </p>
          <ul>
            <li>Supply your name to be used on the scoreboard</li>
            <li>You will need to pick 2 numbers from 1 to 10 and click play</li>
            <li>When at least 2 players choose to play, the game will start</li>
            <li>The game lasts for 30 rounds, each round takes 5 seconds</li>
            <li>You can supply the numbers you want per round or use the previous numbers if you want</li>
            <li>The game pool will generate random numbers and compare with yours and added to your total scores</li>
            <li>If you want to play in the middle of a game, you'll have to wait till the game ends</li>
            <li>Game restarts in 10 seconds</li>
          </ul>

          <p>Scoring System: </p>
          <p>If a player scores exactly 21 points, the game will end with that player</p>
          <p>The player with the highest total score wins, if there is a tie, the upper bound number will be used or the lower bound and their names in alphabetical order to decide the winner. </p>
        </div>
      ),
      onOk() {},
    });
  }


  return (
    <App.Wrapper className="App">
      <h2><strong>Number Pool Game</strong></h2>
      <h4><button className="instruction-button" onClick={gameInstruction}>How To Play</button></h4>
      <div className="game-container">
        <section>
          <CustomAlert info={info} visible={visible} handleClose={handleClose} />
          <CustomAlert info={roundInfo} visible={visibleRoundInfo} handleClose={handleClose} />
        </section>
        <section>
          <p>Play or Watch the game</p>
          <p>
            <Badge status={onlineStatus ? "success" : "default"} text={onlineStatus ? "online" : "offline"} />
            <br />
            {onlineStatus && gameInSession && <Tag color="green">Game In Session</Tag>}
          </p>
          <div className="button-group">
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
            <Button onClick={() => handlePlayerMode("observe")}>Watch</Button></div>
          <div className="scoreboard">{playerMode && <ScoreBoard dataSource={scoreBoard} />}</div>
        </section>
        <section>
          {playerMode === "play" && <PlayerInput onFinish={onFinish} />}
        </section>
      </div>
    </App.Wrapper>
  );
}

App.Wrapper = styled.div`
  padding: 3em;

  .game-container {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    column-gap: 1em;
  }

  .instruction-button {
    border: none;
    background: transparent;
    color: deepskyblue;
    cursor: pointer;

    &:focus {
      outline: none;
    }
  }

  .button-group {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    column-gap: 1em;
  }

  .scoreboard {
    margin-top: 3em;
  }
`;

export default App;
