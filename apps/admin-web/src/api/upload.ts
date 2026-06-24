import client from "./client";

export function uploadFile(file: File) {
  const formData = new FormData();
  formData.append("file", file);
  return client.post("/uploads", formData, {
    headers: { "Content-Type": "multipart/form-data" },
  });
}

export function getUploads(params?: Record<string, any>) {
  return client.get("/uploads", { params });
}

export function deleteUpload(id: number) {
  return client.delete(`/uploads/${id}`);
}
