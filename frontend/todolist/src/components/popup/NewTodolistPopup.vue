<template>
    <MyDialog
      header="Créer une Todolist"
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
          <label for="nom" class="col-12 mb-2">
            Nom
          </label>
          <div class="col-12">
            <InputText
              id="nom"
              placeholder="Nom de la nouvelle Todolist"
              autoResize
              :rows="3"
              v-model="nom"
              v-on:keyup.enter="createTodolist"
            />
          </div>
          <div v-if="error" class="col-12 text-red-500 mt-2">
            {{ error }}
          </div>
        </div>
      </div>
      <template #footer>
        <MyButton label="Annuler" @click="closePopup" severity="secondary" />
        <MyButton
          label="Créer la Todolist"
          :loading="loading"
          @click="createTodolist"
        />
      </template>
    </MyDialog>
  </template>
  
  <script>
  import axios from 'axios';

  export default {
    name: "NewTodolistPopup",
    emits: ["close"],
    data() {
      return {
        nom: "",
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
      createTodolist() {
        if (this.nom.trim() === "") {
          this.error = "Le nom ne peut pas être vide.";
          return;
        }
        this.error = null;
        this.loading = true;
        axios.post('http://localhost:8080/api/v1/todolist',
          {
            nom: this.nom
          }, {
            headers: {
                'Authorization': `Bearer ${sessionStorage.getItem('jwtToken')}`
              }
          }).then((response) => {
            this.todolists = response.data.data;
            this.loading = false;
            this.closePopup(true);
        })
        .catch((error) => {
            console.log(error);
        });
      },
    },
  };
  </script>
  