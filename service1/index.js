// create simple express server
const express = require("express");
const app = express();
const axios = require("axios");
const port = 8199;

app.get("/", async (req, res) => {
  try {
    const response = await axios.get(
      "http://golang-service2:3001/get-container-info",
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
  } catch (error) {
    console.error(
      "Error fetching all container info from golang-service:",
      error.message
    );
    res.status(500).send("Error fetching container info.");
  }
});

app.listen(port, () => {
  console.log(`Node-service listening at http://localhost:${port}`);
});
