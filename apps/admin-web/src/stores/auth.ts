import { defineStore } from "pinia";
import { ref } from "vue";
import router from "@/router";

export const useAuthStore = defineStore("auth", () => {
  const token = ref(localStorage.getItem("admin_token") || "");
  const username = ref(localStorage.getItem("admin_username") || "");
  const nickname = ref(localStorage.getItem("admin_nickname") || "");

  function setAuth(t: string, u: string, n: string) {
    token.value = t;
    username.value = u;
    nickname.value = n;
    localStorage.setItem("admin_token", t);
    localStorage.setItem("admin_username", u);
    localStorage.setItem("admin_nickname", n);
  }

  function logout() {
    token.value = "";
    username.value = "";
    nickname.value = "";
    localStorage.removeItem("admin_token");
    localStorage.removeItem("admin_username");
    localStorage.removeItem("admin_nickname");
    router.push("/login");
  }

  return { token, username, nickname, setAuth, logout };
});
