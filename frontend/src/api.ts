const API_BASE_URL = "http://localhost:8080"; // Assuming the Go server is running on port 8080
const API_KEY = "testkey";

export const uploadFile = async (file: File) => {
  const formData = new FormData();
  formData.append("file", file);

  const response = await fetch(`${API_BASE_URL}/upload`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${API_KEY}`,
    },
    body: formData,
  });

  if (!response.ok) {
    throw new Error("File upload failed");
  }

  return response.json();
};

export const getFiles = async () => {
  const response = await fetch(`${API_BASE_URL}/files`, {
    headers: {
      Authorization: `Bearer ${API_KEY}`,
    },
  });

  if (!response.ok) {
    throw new Error("Failed to fetch files");
  }

  const data = await response.json();
  return data ?? []; // fallback to [] if null
};


export const getFileUrl = (id: string) => {
  return `${API_BASE_URL}/files/${id}`;
};

export const deleteFile = async (id: string) => {
  const response = await fetch(`${API_BASE_URL}/delete/${id}`, {
    method: "DELETE",
    headers: {
      Authorization: `Bearer ${API_KEY}`,
    },
  });

  if (!response.ok) {
    throw new Error("Failed to delete file");
  }

  return response.json();
};

export const updateFileName = async (id: string, newName: string) => {
  const response = await fetch(`/update/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ name: newName }),
  });

  if (!response.ok) {
    throw new Error('Failed to update file name');
  }

  return response.json();
};
