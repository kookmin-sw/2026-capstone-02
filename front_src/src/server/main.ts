
import express from "express";
import ViteExpress from "vite-express";

import multer from "multer";

import fs from "fs";
import path from "path";
import url from "url";

const app = express();
const port = 3000;

const __filename = url.fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Multer storage to save files to a specific folder
const UPLOAD_DIR = path.resolve(__dirname, "../../go");

if (!fs.existsSync(UPLOAD_DIR))
    fs.mkdirSync(UPLOAD_DIR, { recursive: true });

const storage = multer.diskStorage({
    destination: (_req, _file, cb) => cb(null, UPLOAD_DIR),
    filename: (_req, file, cb) => cb(null, file.originalname),
});

const upload = multer({ storage });

// Single-file upload endpoint
app.post("/upload", upload.single("file"), (req, res) => {
    if (!req.file)
        return res.status(400).json({ error: "No file uploaded" });

    const savedName = req.file.filename;
    const savedPath = req.file.path;

    if (!savedName.toLowerCase().endsWith(".go")) {
        fs.unlink(savedPath, () => { });
        return res.status(400).json({ error: "Only .go files are allowed" });
    }

    res.json({ savedAs: savedName, path: savedPath });
});

ViteExpress.listen(app, port, () =>
    console.log(`Server is listening on ${port}...`),
);
