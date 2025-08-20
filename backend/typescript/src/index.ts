import express from "express";
import multer from "multer";
import dotenv from "dotenv";
import {
  uploadHandler,
  listHandler,
  downloadHandler,
  deleteHandler,
} from "./handlers";

dotenv.config();

const app = express();
const PORT = process.env.PORT || 8080;
const API_KEY = process.env.API_KEY;

const upload = multer({ dest: "uploads/" });

// Middleware for API key
function requireAPIKey(req: express.Request, res: express.Response, next: express.NextFunction) {
  const auth = req.headers.authorization;
  if (!auth || auth !== `Bearer ${API_KEY}`) return res.status(401).send("Unauthorized");
  next();
}

// Root instructions
app.get("/", (req, res) => {
  res.send(`
File Server API

Endpoints:
POST   /upload        (auth required)
GET    /files         (auth required)
GET    /files/:id     (public)
DELETE /delete/:id    (auth required)

Auth: Authorization: Bearer <API_KEY>
`);
});

// Routes
app.post("/upload", requireAPIKey, upload.single("file"), uploadHandler);
app.get("/files", requireAPIKey, listHandler);
app.get("/files/:id", downloadHandler);
app.delete("/delete/:id", requireAPIKey, deleteHandler);

app.listen(PORT, () => {
  console.log(`File server running on http://localhost:${PORT}`);
});
