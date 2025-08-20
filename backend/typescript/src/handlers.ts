import { Request, Response } from "express";
import fs from "fs";
import path from "path";
import { generateFileId } from "./utils";

const dataDir = path.join(__dirname, "../uploads");
if (!fs.existsSync(dataDir)) fs.mkdirSync(dataDir);

interface FileMeta {
  id: string;
  filename: string;
  originalName: string;
}

export const store: { [id: string]: FileMeta } = {};

// Upload handler
export function uploadHandler(req: Request, res: Response) {
  if (!req.file) return res.status(400).send("No file uploaded");

  const id = generateFileId(req.file.originalname);
  const ext = path.extname(req.file.originalname);
  const storedFileName = `${id}${ext}`;

  fs.renameSync(req.file.path, path.join(dataDir, storedFileName));

  store[id] = {
    id,
    filename: storedFileName,
    originalName: req.file.originalname,
  };

  res.send(id);
}

// List files
export function listHandler(req: Request, res: Response) {
  res.json(Object.values(store));
}

// Download/View file
export function downloadHandler(req: Request, res: Response) {
  const id = req.params.id;

  if (!id) return res.status(400).send("Missing file id");

  const fileMeta = store[id];
  if (!fileMeta) return res.status(404).send("File not found");

  res.sendFile(path.join(dataDir, fileMeta.filename));
}

// Delete file
export function deleteHandler(req: Request, res: Response) {
  const id = req.params.id;

  if (!id) return res.status(400).send("Missing file id");

  const fileMeta = store[id];
  
  if (!fileMeta) return res.status(404).send("File not found");

  fs.unlinkSync(path.join(dataDir, fileMeta.filename));
  delete store[id];
  res.send("Deleted");
}
