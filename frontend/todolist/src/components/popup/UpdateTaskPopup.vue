<template>
    <MyDialog
      header="Créer une Tâche"
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
              placeholder="Nom de la tache"
              autoResize
              :rows="3"
              v-model="nom"
            />
          </div>
          <div v-if="error" class="col-12 text-red-500 mt-2">
            {{ error }}
          </div>
        </div>
        <div class="field grid">
          <label for="description" class="col-12 mb-2">
            Description
          </label>
          <div class="col-12">
            <TextArea
              id="description"
              placeholder="Description de la tache"
              autoResize
              :rows="3"
              v-model="description"
            />
          </div>
        </div>
        <div class="field grid">
          <label for="date" class="col-12 mb-2">
            Date d'échéance
          </label>
          <div class="col-12">
            <MyCalendar v-model="date" dateFormat="dd/mm/yy" showIcon :manualInput="false"></MyCalendar>
          </div>
        </div>
        <div class="field grid">
          <label for="Priorité" class="col-12 mb-2">
            Priorité
          </label>
          <div class="col-12">
            <div class="flex flex-wrap gap-3">
              <div v-for="option in options" :key="option.id" class="flex align-items-center">
                <MyRadioButton v-model="priorité" :inputId="option.id" name="dynamic" :value="option.value" />
                <label :for="option.id" class="ml-2">{{ option.label }}</label>
              </div>
            </div>
          </div>
        </div>
        <div>
    </div>
      </div>
      <template #footer>
        <MyButton label="Annuler" @click="closePopup" severity="secondary" />
        <MyButton
          label="Mettre à jour la tâche"
          :loading="loading"
          @click="updateTask"
        />
      </template>
    </MyDialog>
  </template>
  
  <script>
  import axios from 'axios';
    export default {
        name: "UpdateTask",
        emits: ["close"],
        props: {
            task: null,
        },
        data() {
            return {
                nom: "",
                description: "",
                priorité: 0,
                loading: false,
                date: null,
                options: [
                  { id: 'option1', value: 0, label: 'Faible' },
                  { id: 'option2', value: 1, label: 'Moyenne' },
                  { id: 'option3', value: 2, label: 'Forte' },
                  { id: 'option4', value: 3, label: 'Urgente' },
                ],
                error: null, 
            };
        },
        mounted() {
            this.nom = this.task.nom;
            this.description = this.task.description;
            this.priorité = this.task.priorite;

            const parts = this.task.deadline.split("/");
            const formattedDate = `${parts[2]}-${parts[1]}-${parts[0]}`;
            const dateObject = new Date(formattedDate);
            this.date = dateObject;
        },
        methods: {
            closePopup(reload) {
                this.$emit("close", reload);
            },
            formatDate(date) {
                const day = date.getDate();
                const month = date.getMonth() + 1;
                const year = date.getFullYear();

                const formattedDay = (day < 10) ? `0${day}` : day;
                const formattedMonth = (month < 10) ? `0${month}` : month;
                
                const formattedDate = `${formattedDay}/${formattedMonth}/${year}`;

                return formattedDate;
            },
            updateTask() {
              if (this.nom.trim() === "") {
                this.error = "Le nom ne peut pas être vide.";
                return;
              }
              const item = {
                  "nom": this.nom,
                  "description": this.description,
                  "priorite": parseInt(this.priorité),
                  "etat": this.task.etat,
                  "deadline": this.formatDate(this.date),
              }
              const url = 'http://localhost:8080/api/v1/task/' + this.task.idTache;
              axios.put(url, item, {
                headers: {
                    'Authorization': `Bearer ${sessionStorage.getItem('jwtToken')}`
                  }
              })
              .then(()=> {
                  this.$emit("close", true);
              })
              .catch((error) => {
                  console.log(error);
              });
            }
        },
    };
  </script>
  