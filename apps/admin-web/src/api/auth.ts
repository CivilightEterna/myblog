import client from "./client";

export function login(username: string, password: string) {
  return client.post("/login", { username, password });
}

export function getProfile() {
  return client.get("/profile");
}
