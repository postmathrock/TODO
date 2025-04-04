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
  const [editModalIsOpen, setEditModalIsOpen] = useState<boolean>(false);

  const [selectedTask, setSelectedTask] = useState<Task>({
    "id": -1,
    "body": "",
    "time_limit": "",
    "completed_at": "",
    "created_at": ""
  });

  // タスク作成ボタンのクリックハンドラー
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

  // タスク更新ボタンのクリックハンドラー
  function clickUpdateButton(id: number) {
    fetch("http://localhost:8080/api/tasks/update", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ "ID": id, "body": body, "time_limit": time_limits })
    }).then(() => getTasks())
  }

  // タスク削除ボタンのクリックハンドラー
  function clickDeleteButton(id: number) {
    fetch("http://localhost:8080/api/tasks/delete", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ "ID": id})
    }).then(() => getTasks())
  }

  // タスク完了ボタンのクリックハンドラー
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
      {/* アプリのタイトル */}
      <h1>TODO app</h1>

      {/* タスクの入力欄 */}
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

            {/* 編集アイコン */}
            <EditIcon onClick={() => {
              setSelectedTask(val);
              setEditModalIsOpen(true);
            }} />
            {/* 編集のモーダル */}
            <Modal isOpen={editModalIsOpen}>
              
              {/* モーダルの閉じるボタン */}
              <Button onClick={() => {
                setEditModalIsOpen(false);
              }} >X</Button>

              {/* 編集するタスクの入力欄 */}
              <TextField id="outlined-basic" label="Outlined" variant="outlined" defaultValue={selectedTask.body} onChange={onChangeBody} />

              {/* 期限日の編集 */}
              <LocalizationProvider dateAdapter={AdapterDateFns}>
                <DatePicker label="Basic date picker" onChange={onChangeTimeLimit} />
              </LocalizationProvider>

              {/* 更新ボタン */}
              {/* Buttonの表示 */}
              <Button onClick={() => {
                setEditModalIsOpen(false);
                {/* タスクの更新 */}
                clickUpdateButton(selectedTask.id)
              }}>update</Button>
            </Modal>

            {/* 削除アイコン */}
            {/* DeleteIconの表示 */}
            <DeleteIcon onClick={() => {
            {/*  */}
            setSelectedTask(val)
            {/* タスクの削除 */}
            clickDeleteButton(selectedTask.id)
            }} />

            {/* タスク完了のチェックボックス */}
            {/* 値が入っていなければCheckBoxOutlineBlankIconをレンダリング(表示)し、クリックすると完了状態に切り替える為の関数を呼ぶ */}
            {!val.completed_at && <CheckBoxOutlineBlankIcon onClick={() => {clickCheckBox(val.id)}} />}
            {/* 値が入っていればCheckBoxIconをレンダリング(表示)する */}
            {val.completed_at && <CheckBoxIcon />}
          </Card>
        )
      }
    </>

  )
}

export default App
