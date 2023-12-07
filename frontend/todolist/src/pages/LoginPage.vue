<template>
    <div class="flex items-center justify-content-center mt-8">
        <div class="surface-card p-4 shadow-2 border-round w-full lg:w-6">
            <div class="text-center mb-5">
                <div class="text-900 text-3xl font-medium mb-3">Connexion</div>
                <span class="text-600 font-medium line-height-3">Vous n'avez pas de compte ?</span>
                <a class="font-medium no-underline ml-2 text-blue-500 cursor-pointer" @click="CreateAccount">Créer un compte</a>
            </div>
            <div v-if="error" class="col-12 text-red-500 mt-2">
                {{ error }}
            </div>
            <div>
                <label for="username" class="block text-900 font-medium mb-2">Nom d'utilisateur</label>
                <InputText id="username" v-model="username" type="text" class="w-full mb-3" />
                <label for="password" class="block text-900 font-medium mb-2">Mot de passe</label>
                <InputText id="password" type="password" v-model="password" v-on:keyup.enter="Connect" class="w-full mb-3" />
                <!-- 
                    <div class="flex align-items-center justify-content-between mb-6">
                        <a class="font-medium no-underline ml-2 text-blue-500 text-right cursor-pointer">Forgot password?</a>
                    </div> 
                -->
                <MyButton label="Se connecter" icon="pi pi-user" class="w-full" @click="Connect"></MyButton>
            </div>
        </div>
  </div>
  <CreateAccountPopup v-if="popup === 'CREATEACCOUNT'" @close="closePopup" />
</template>
  
<script>
import router from '@/router';
import axios from 'axios';
import CreateAccountPopup from '@/components/popup/CreateAccountPopup.vue';

  export default {
    name: 'LoginPage',
    components: { CreateAccountPopup },
    created() {
        if (sessionStorage.getItem('jwtToken')) {
        router.push({ path: '/calendar' });
        }
    },
    data() {
        return {
            error: null,
            username: "",
            password: "",
            popup: null,
        };
    },
    methods: {
        closePopup() {
            this.popup = null;
        },
        CreateAccount() {
            this.popup = "CREATEACCOUNT";
        },
        Connect() {
            if (this.username.trim() === "" || this.password.trim() === "") {
                this.error = "Le nom d'utilisateur ou le mot de passe ne peuvent pas être vides.";
                return;
            }
            const item = {
                username: this.username,
                password: this.password,
            }
            const url = 'http://localhost:8080/login';
            axios.post(url, item)
            .then((response)=> {
                sessionStorage.setItem('jwtToken', response.data.token);
                sessionStorage.setItem('username', response.data.username);
                router.push({path: '/calendar'});
            })
            .catch(() => {
                this.error = "Nom d'utilisateur ou mot de passe incorrects."
                this.password = "";
            });
        }
    },
  };
  </script>