<template>
  <div class="card">
    <MyToolbar
      v-if="isUserLoggedIn"
      class="bg-gray-900 shadow-2"
      style="border-radius: 15px; background-image: linear-gradient(to right, var(--bluegray-500), var(--bluegray-800))"
    >
      <template #start> <span class="font-bold text-bluegray-50">{{ username }}</span> </template>
      <template #center>
        <div class="flex items-center justify-center gap-3 flex-shrink-0">
          <router-link to="/urgent" class="p-link inline-flex justify-content-center align-items-center text-white h-3rem w-3rem border-circle hover:bg-white-alpha-10 transition-all transition-duration-200" style="text-decoration: none;">
            <i class="pi pi-exclamation-triangle text-2xl"></i>
          </router-link>
          <router-link to="/calendar" class="p-link inline-flex justify-content-center align-items-center text-white h-3rem w-3rem border-circle hover:bg-white-alpha-10 transition-all transition-duration-200" style="text-decoration: none;">
            <i class="pi pi-calendar text-2xl"></i>
          </router-link>
          <router-link to="/list" class="p-link inline-flex justify-content-center align-items-center text-white h-3rem w-3rem border-circle hover:bg-white-alpha-10 transition-all transition-duration-200" style="text-decoration: none;">
            <i class="pi pi-list text-2xl"></i>
          </router-link>
          <MyButton link class="p-link inline-flex justify-content-center align-items-center text-white h-3rem w-3rem border-circle hover:bg-white-alpha-10 transition-all transition-duration-200" @click="toggleTheme" :icon="themeIcon"></MyButton>
        </div>
      </template>
      <template #end>
        <MyButton @click="disconnect" outlined rounded class="ml-8" icon="ri ri-logout-circle-line"></MyButton>
      </template>
    </MyToolbar>
    <MyToolbar
      v-else
      class="bg-gray-900 shadow-2"
      style="border-radius: 15px; background-image: linear-gradient(to right, var(--bluegray-500), var(--bluegray-800))"
    >
      <template #center>
        <div class="flex items-center justify-center gap-3 flex-shrink-0">
          <MyButton link class="p-link inline-flex justify-content-center align-items-center text-white h-3rem w-3rem border-circle hover:bg-white-alpha-10 transition-all transition-duration-200" @click="toggleTheme" :icon="themeIcon"></MyButton>
        </div>
      </template>
    </MyToolbar>
  </div>
</template>

<script>
import router from '@/router';
  export default {
    name: 'TopBar',
    data() {
      return {
        isUserLoggedIn: sessionStorage.getItem('jwtToken') !== null,
        currentTheme: localStorage.getItem('theme') || 'lara-light-cyan',
        username: sessionStorage.getItem('username') || '',
      };
    },
    computed: {
      themeIcon() {
        return this.currentTheme === 'lara-dark-cyan' ? 'pi pi-sun' : 'pi pi-moon';
      }
    },
    created() {
      this.initTheme();
    },
    methods: {
      disconnect() {
        sessionStorage.removeItem('jwtToken');
        sessionStorage.removeItem('username');
        router.push({ name: 'login' });
      },
      initTheme() {
        const head = document.head; 
        const link = document.getElementById('theme-link');

        const newLink = document.createElement('link');
        newLink.id = 'theme-link';
        newLink.rel = 'stylesheet';
        newLink.href = `https://unpkg.com/primevue/resources/themes/${this.currentTheme}/theme.css`;
        if (link) {
          head.removeChild(link);
        }
        head.appendChild(newLink);
        localStorage.setItem('theme', this.currentTheme);
      },
      toggleTheme() {
        let nextTheme = '';
        if (this.currentTheme === 'lara-light-cyan') nextTheme = 'lara-dark-cyan';
        else if (this.currentTheme === 'lara-dark-cyan') nextTheme = 'lara-light-cyan';

        const head = document.head;
        const link = document.getElementById('theme-link');

        const newLink = document.createElement('link');
        newLink.id = 'theme-link';
        newLink.rel = 'stylesheet';
        newLink.href = `https://unpkg.com/primevue/resources/themes/${nextTheme}/theme.css`;

        if (link) {
          head.removeChild(link);
        }

        head.appendChild(newLink);
        this.currentTheme = nextTheme;

        localStorage.setItem('theme', this.currentTheme);
      },
    },
    watch: {
      $route() {
        this.isUserLoggedIn = sessionStorage.getItem('jwtToken') !== null;
        this.username = sessionStorage.getItem('username') || '';
      }
    },
  };
</script>

<style scoped></style>
