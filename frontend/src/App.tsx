import { useState, useEffect } from 'react';
import { uploadFile, getFiles, getFileUrl, deleteFile } from './api';
import './App.css';

interface File {
  id: string;
  name: string;
  size: number;
  created_at: string;
}

function App() {
  const [files, setFiles] = useState<File[]>([]);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);

  useEffect(() => {
    fetchFiles();
  }, []);

  const fetchFiles = async () => {
    try {
      const files = await getFiles();
      console.log(files);
      setFiles(files);
    } catch (error) {
      console.error(error);
    }
  };

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.files) {
      setSelectedFile(event.target.files[0]);
    }
  };

  const handleUpload = async () => {
    if (selectedFile) {
      try {
        await uploadFile(selectedFile);
        fetchFiles();
        setSelectedFile(null);
      } catch (error) {
        console.error(error);
      }
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await deleteFile(id);
      fetchFiles();
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>File Server</h1>
      </header>
      <main>
        <div className="upload-section">
          <input type="file" onChange={handleFileChange} />
          <button onClick={handleUpload} disabled={!selectedFile}>
            Upload
          </button>
        </div>
        <div className="file-list">
          <h2>Uploaded Files</h2>
          <ul>
            {files.map((file) => (
              <li key={file.id}>
                <a href={getFileUrl(file.id)} target="_blank" rel="noopener noreferrer">
                  {file.name}
                </a>
                <span>({(file.size / 1024).toFixed(2)} KB)</span>
                <button onClick={() => handleDelete(file.id)}>Delete</button>
              </li>
            ))}
          </ul>
        </div>
      </main>
    </div>
  );
}

export default App;