import { useState } from "react";

const ROOT_URL = "/api";
// 'http://localhost:3000/api/containers'

function App() {
  const [text, setText] = useState("");

  const fetchContainerInfo = async () => {
    const response = await fetch(`${ROOT_URL}/get-container-info`);

    const text = await response.text();

    setText(text);
  };

  const stopAllContainers = async () => {
    const response = await fetch(`${ROOT_URL}/stop-all-containers`, {
      method: "POST",
    });

    if (response.status === 200) {
      console.log("All containers stopped");
    } else {
      console.error("Error stopping all containers");
    }
  };

  return (
    <div>
      <p>Get all containers</p>
      <button onClick={() => fetchContainerInfo()}>REQUEST</button>
      <br />
      <textarea
        style={{ display: "block" }}
        cols={100}
        rows={20}
        value={text}
        readOnly
      ></textarea>

      <br />

      <button onClick={() => stopAllContainers()}>STOP</button>
    </div>
  );
}

export default App;
