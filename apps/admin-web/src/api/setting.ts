import client from "./client";

export function getSettings() {
  return client.get("/settings");
}

export function updateSettings(settings: Record<string, string>) {
  return client.put("/settings", settings);
}
