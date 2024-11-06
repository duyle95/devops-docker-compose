// create simple express server
const express = require("express");
const app = express();
const axios = require("axios");
const fs = require("fs");
const path = require("path");
const cors = require("cors");
const port = 3000;

let isSleeping = false;

app.use(cors());

app.use((req, res, next) => {
  if (isSleeping) {
    return res
      .status(503)
      .send("Server is temporarily unavailable, please try again later.");
  }
  next();
});

app.get("/api/get-container-info", async (req, res) => {
  console.log("GET /api/get-container-info");
  try {
    const response = await axios.get(
      "http://golang-service:3001/get-container-info",
      {
        responseType: "text",
      }
    );

    res.setHeader("Content-Type", "text/plain");
    res.setHeader(
      "Content-Disposition",
      "attachment; filename=duyle-container-info.txt"
    );
    res.send(response.data);

    isSleeping = true;
    await new Promise((resolve) =>
      setTimeout(() => {
        isSleeping = false;
        resolve();
      }, 2000)
    );
  } catch (error) {
    console.error(
      "Error fetching all container info from golang-service:",
      error.message
    );
    res.status(500).send("Error fetching container info.");
  }
});

app.post("/api/stop-all-containers", async (req, res) => {
  console.log("POST /api/stop-all-containers");
  try {
    const response = await axios.post(
      "http://golang-service:3001/stop-all-containers"
    );

    res.sendStatus(200);
  } catch (error) {
    console.error("Error when stopping all containers: ", error.message);
    res.status(500);
  }
});

app.listen(port, () => {
  console.log(`Node-service listening at port ${port}`);
});
