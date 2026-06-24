import client from "./client";

export function triggerBuild(triggerType = "manual") {
  return client.post("/build", { trigger_type: triggerType });
}

export function getBuilds(params?: Record<string, any>) {
  return client.get("/builds", { params });
}

export function getBuild(id: number) {
  return client.get(`/builds/${id}`);
}
