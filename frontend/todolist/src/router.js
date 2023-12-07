import { createRouter, createWebHistory } from "vue-router";
import HomePage from "./pages/HomePage.vue";
import ListPage from "./pages/ListPage.vue";
import UrgentPage from "./pages/UrgentPage.vue";
import LoginPage from "./pages/LoginPage.vue";

const routes = [
  {
    path: "/",
    redirect: "/login",
  },
  {
    path: "/login",
    name: "login",
    meta: {
        requiresAuth: false,
    },
    component: LoginPage,
  },
  {
    path: "/calendar",
    name: "home",
    meta: {
        requiresAuth: true,
    },
    component: HomePage,
  },
  {
    path: "/list",
    name: "listPage",
    meta: {
        requiresAuth: true,
    },
    component: ListPage,
  },
  {
    path: "/urgent",
    name: "urgentPage",
    meta: {
        requiresAuth: true,
    },
    component: UrgentPage,
  }
];

const router = createRouter({
  history: createWebHistory("/"),
  routes,
});


router.beforeEach((to, from, next) => {
  if (to.meta.requiresAuth) {
    const token = sessionStorage.getItem('jwtToken');
    if (!token) {
      next({ name: 'login' });
    } else {
      next();
    }
  } else {
    next();
  }
});

export default router;