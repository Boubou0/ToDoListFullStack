<template>
    <MyDialog
      header="Créer un compte"
      visible
      modal
      @keyup.esc="closePopup(false)"
      class="popup"
    >
      <template #closeicon>
        <MyButton
          icon="ri-close-line"
          @click="closePopup(false)"
          outlined
          text
          severity="secondary"
          autofocus="false"
        />
      </template>
      <div class="p-fluid">
        <div class="field grid">
          <label for="username" class="col-12 mb-2">
            Nom d'utilisateur
          </label>
          <div class="col-12">
            <InputText
              id="username"
              placeholder="Nom d'utilisateur"
              autoResize
              :rows="3"
              v-model="username"
            />
          </div>
        </div>
        <div class="field grid">
          <label for="password" class="col-12 mb-2">
            Mot de passe
          </label>
          <div class="col-12">
            <InputText
              id="password"
              type="password"
              placeholder="Mot de passe"
              autoResize
              :rows="3"
              v-model="password"
              v-on:keyup.enter="createAccount"
            />
          </div>
        </div>
        <div class="field grid">
          <label for="password2" class="col-12 mb-2">
            Confirmer le mot de passe
          </label>
          <div class="col-12">
            <InputText
              id="password2"
              type="password"
              placeholder="Mot de passe"
              autoResize
              :rows="3"
              v-model="password2"
              v-on:keyup.enter="createAccount"
            />
          </div>
        </div>
      </div>
      <div v-if="error" class="col-12 text-red-500 mt-2">
        {{ error }}
        </div>
      <template #footer>
        <MyButton label="Annuler" @click="closePopup" severity="secondary" />
        <MyButton
          label="Créer le compte"
          :loading="loading"
          @click="createAccount"
        />
      </template>
    </MyDialog>
  </template>
  
  <script>
import axios from 'axios';

  export default {
    name: "CreateAccountPopup",
    emits: ["close"],
    data() {
      return {
        username: "",
        password: "",
        password2: "",
        loading: false,
        error: null,
      };
    },
    mounted() {
      this.loading = false;
    },
    methods: {
      closePopup(reload) {
        this.$emit("close", reload);
      },
      createAccount() {
        if (this.username.trim() === "") {
          this.error = "Le nom d'utilisateur ou le mot de passe ne peuvent pas être vides.";
          return;
        }
        if (this.password.trim() === "") {
          this.error = "Le nom d'utilisateur ou le mot de passe ne peuvent pas être vides.";
          return;
        }
        if (this.password != this.password2){
          this.error = "La confirmation du mot de passe doit correspondre au mot de passe saisi précédemment."
          this.password = "";
          this.password2 = "";
          return;
        }
        const item = {
            "username": this.username,
            "password": this.password,
        }
        const url = 'http://localhost:8080/signup';
            axios.post(url, item)
            .then(()=> {
                this.closePopup(false);
            })
            .catch((error) => {
                if(error.response.status === 409) {
                    this.error = "Ce nom d'utilisateur est déjà utilisé"
                } else {
                    this.error = error.response.data.error;
                }
            });
      }
    },
  };
  </script>
  