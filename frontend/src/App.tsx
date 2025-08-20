import { useState, useEffect } from 'react';
import { uploadFile, getFiles, getFileUrl, deleteFile,updateFileName } from './api';
import './App.css';

interface File {
  id: string;
  name: string;
  size: number;
  created_at: string;
}

function App() {
  const [files, setFiles] = useState<File[]>([]);
  const [isUploading, setIsUploading] = useState(false);
  const [dragActive, setDragActive] = useState(false);

  useEffect(() => {
    fetchFiles();
    
    // Add page-wide drag and drop listeners
    const handlePageDragOver = (e: DragEvent) => {
      e.preventDefault();
      setDragActive(true);
    };

    const handlePageDragLeave = (e: DragEvent) => {
      // Only hide drag indicator when leaving the window entirely
      if (!e.relatedTarget) {
        setDragActive(false);
      }
    };

    const handlePageDrop = async (e: DragEvent) => {
      e.preventDefault();
      setDragActive(false);
      
      if (e.dataTransfer?.files && e.dataTransfer.files[0]) {
        const file = e.dataTransfer.files[0];
        await uploadFileDirectly(file);
      }
    };




    document.addEventListener('dragover', handlePageDragOver);
    document.addEventListener('dragleave', handlePageDragLeave);
    document.addEventListener('drop', handlePageDrop);

    return () => {
      document.removeEventListener('dragover', handlePageDragOver);
      document.removeEventListener('dragleave', handlePageDragLeave);
      document.removeEventListener('drop', handlePageDrop);
    };
  }, []);

  const fetchFiles = async () => {
    try {
      const files = await getFiles();
      setFiles(files);
    } catch (error) {
      console.error(error);
    }
  };

  const uploadFileDirectly = async (file: File) => {
    setIsUploading(true);
    try {
      await uploadFile(file);
      await fetchFiles();
      // Reset file input
      const fileInput = document.getElementById('file-input') as HTMLInputElement;
      if (fileInput) fileInput.value = '';
    } catch (error) {
      console.error(error);
    } finally {
      setIsUploading(false);
    }
  };

  const handleFileChange = async (event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.files && event.target.files[0]) {
      const file = event.target.files[0];
      await uploadFileDirectly(file);
    }
  };

  const handleDelete = async (id: string, fileName: string) => {
    if (window.confirm(`Are you sure you want to delete "${fileName}"?`)) {
      try {
        await deleteFile(id);
        await fetchFiles();
      } catch (error) {
        console.error(error);
      }
    }
  };

      const handleRename = async (id: string, currentName: string) => {
  const newName = window.prompt('Enter new file name:', currentName);
  if (!newName || newName === currentName) return;

  try {
    await updateFileName(id, newName);
    await fetchFiles();
  } catch (error) {
    console.error(error);
    alert('Failed to update file name.');
  }
};

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const getFileIcon = (fileName: string) => {
    const extension = fileName.split('.').pop()?.toLowerCase();
    switch (extension) {
      case 'pdf': return '📄';
      case 'doc':
      case 'docx': return '📝';
      case 'jpg':
      case 'jpeg':
      case 'png':
      case 'gif': return '🖼️';
      case 'mp4':
      case 'avi':
      case 'mov': return '🎥';
      case 'mp3':
      case 'wav': return '🎵';
      case 'zip':
      case 'rar': return '📦';
      default: return '📁';
    }
  };

  return (
    <div className="app">
      {/* Page-wide drag overlay */}
      {dragActive && (
        <div className="drag-overlay">
          <div className="drag-message">
            <div className="drag-icon">📁</div>
            <h2>Drop files anywhere to upload</h2>
            <p>Release to upload instantly</p>
          </div>
        </div>
      )}

      <header className="header">
        <div className="header-content">
          <h1>📂 File Manager</h1>
          <p>Upload, manage, and share your files</p>
        </div>
      </header>

      <main className="main-content">
        <div className="upload-container">
          <div className="upload-area">
            <div className="upload-icon">☁️</div>
            <h3>Select a file to upload instantly</h3>
            <p>Click the button below or drag files anywhere on the page</p>
            
            <input
              id="file-input"
              type="file"
              onChange={handleFileChange}
              className="file-input"
              disabled={isUploading}
            />
            
            <label 
              htmlFor="file-input" 
              className={`browse-btn ${isUploading ? 'uploading' : ''}`}
            >
              {isUploading ? (
                <>
                  <span className="spinner"></span>
                  Uploading...
                </>
              ) : (
                <>
                  📁 Choose File
                </>
              )}
            </label>
          </div>
        </div>

        <div className="files-container">
          <div className="files-header">
            <h2>Your Files</h2>
            <span className="files-count">{files.length} files</span>
          </div>

          {files.length === 0 ? (
            <div className="empty-state">
              <div className="empty-icon">📁</div>
              <h3>No files uploaded yet</h3>
              <p>Upload your first file to get started</p>
            </div>
          ) : (
            <div className="files-grid">
              {files.map((file) => (
                <div key={file.id} className="file-card">
                  <div className="file-icon-large">
                    {getFileIcon(file.name)}
                  </div>
                  
                  <div className="file-info">
                    <h4 className="file-name" title={file.name}>
                      {file.name.length > 25 ? `${file.name.substring(0, 25)}...` : file.name}
                    </h4>
                    <p className="file-size">{formatFileSize(file.size)}</p>
                    <p className="file-date">{formatDate(file.created_at)}</p>
                  </div>

                  
                  <div className="file-actions">
  <a
    href={getFileUrl(file.id)}
    target="_blank"
    rel="noopener noreferrer"
    className="action-btn download-btn"
    title="Download"
  >
    ⬇️
  </a>

  <button
    onClick={() => handleRename(file.id, file.name)}
    className="action-btn rename-btn"
    title="Rename"
  >
    ✏️
  </button>

  <button
    onClick={() => handleDelete(file.id, file.name)}
    className="action-btn delete-btn"
    title="Delete"
  >
    🗑️
  </button>
</div>

                </div>
              ))}
            </div>
          )}
        </div>
      </main>
    </div>
  );
}

export default App;