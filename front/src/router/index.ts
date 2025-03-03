import { createRouter, createWebHistory } from "vue-router"
import HomeView from "../views/HomeView.vue"
import LoginView from "../views/auth/LoginView.vue"
import LogoutView from "@/views/auth/LogoutView.vue";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "home",
      component: HomeView,
    },
    {
      path: "/login",
      name: "login",
      component: LoginView
    },
    {
      path: "/logout",
      name: "logout",
      component: LogoutView
    }
  ],
})

export default router
