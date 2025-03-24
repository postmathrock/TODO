import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'
import { Button, Card, TextField } from '@mui/material'
import { LocalizationProvider, DatePicker } from '@mui/x-date-pickers'
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFnsV3'
import CheckBoxOutlineBlankIcon from '@mui/icons-material/CheckBoxOutlineBlank';
import CheckBoxIcon from '@mui/icons-material/CheckBox';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import React, { useEffect } from 'react';
import Modal from "react-modal";


// タスクの型
type Task = {
  "id": number;
  "body": string;
  "time_limit": string;
  "completed_at": string;
  "created_at": string;
};

function App() {
  // タスクのデータの保持
  const [tasks, setTasks] = useState<Task[]>([]);
  // 本文のデータの保持
  const [body, setBody] = useState<string>("");
  const onChangeBody = (event) => {      //入力内容が変化したときに実行する関数
    setBody(event.target.value);
  };
  // 期限日のデータの保持
  const [time_limits, setTime_limit] = useState<string>("");
  const onChangeTimeLimit = (event) => {
    setTime_limit(event);
  }
  const [editModalIsOpen, setEditModalIsOpen] = useState(false);

  const [selectedTask, setSelectedTask] = useState<Task>({
    "id": -1,
    "body": "",
    "time_limit": "",
    "completed_at": "",
    "created_at": ""
  });

  function clickSubmitButton() {
    console.log(body, time_limits);
    // tasksのapiに接続
    fetch("http://localhost:8080/api/tasks/create", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ "body": body, "time_limit": time_limits })
    }).then(() => getTasks())
  }
/*
  function clickEditButton() {
    console.log("クリックしました")
  }*/

  function clickUpdateButton(id: number) {
    fetch("http://localhost:8080/api/tasks/update", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ "ID": id, "body": body, "time_limit": time_limits })
    }).then(() => getTasks())
  }

  function clickDeleteButton(id: number) {
    fetch("http://localhost:8080/api/tasks/delete", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ "ID": id})
    }).then(() => getTasks())
  }

  function clickCheckBox (id: number) {
    fetch("http://localhost:8080/api/tasks/complete", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ "ID": id})
    }).then(() => getTasks())
  }

  // useStateから値を受け取り処理する間数　
  function getTasks () {
          //tasksのapiに接続
          fetch("http://localhost:8080/api/tasks")
          // 接続した結果をjsonに変換する
          .then(response => response.json())
          // tasksにデータをセットする
          .then(data => {
            if(data != null && data != undefined){
              setTasks(data);
            }
          })
          .catch(error => {
            console.error(error);
          });
  }

  // getTasksを受け取る
  useEffect(
    () => {
      getTasks()
    },
    []
  );

  // ページの表示
  return (
    <>
      {/* タイトル */}
      <h1>TODO app</h1>

      {/* タスクの入力画面 */}
      <TextField id="outlined-basic" label="Task name" variant="outlined" onChange={onChangeBody} />

      {/* 日付の設定画面 */}
      <LocalizationProvider dateAdapter={AdapterDateFns}>
        <DatePicker label="Basic date picker" onChange={onChangeTimeLimit} />
      </LocalizationProvider>

      {/* タスクの送信ボタン */}
      <Button variant="outlined" onClick={clickSubmitButton}>SUBMIT</Button>

      {/* カードの表示画面 */}
      {
        tasks.map((val) =>
          <Card variant="outlined">
            {/* 本文の表示 */}
            <p>{val.body}</p>

            {/* 期限日の表示 */}
            <p>{val.time_limit}</p>

            {/* 編集ボタン */}
            <EditIcon onClick={() => {
              setSelectedTask(val);
              setEditModalIsOpen(true);
            }} />
            <Modal isOpen={editModalIsOpen}>

              <Button onClick={() => {
                setEditModalIsOpen(false);
              }} >X</Button>
              <TextField id="outlined-basic" label="Outlined" variant="outlined" defaultValue={selectedTask.body} onChange={onChangeBody} />

              {/* 期限日の編集 */}
              <LocalizationProvider dateAdapter={AdapterDateFns}>
                <DatePicker label="Basic date picker" onChange={onChangeTimeLimit} />
              </LocalizationProvider>

              <Button onClick={() => {
                setEditModalIsOpen(false);
                clickUpdateButton(selectedTask.id)
              }}>update</Button>
            </Modal>

            {/* 削除ボタン */}
            <DeleteIcon onClick={() => {
              setSelectedTask(val)
              clickDeleteButton(selectedTask.id)
            }} />

            {/* チェックボックス */}
            {!val.completed_at && <CheckBoxOutlineBlankIcon onClick={() => {clickCheckBox(val.id)}} />}
            {val.completed_at && <CheckBoxIcon />}
          </Card>
        )
      }
    </>

  )
}

export default App
