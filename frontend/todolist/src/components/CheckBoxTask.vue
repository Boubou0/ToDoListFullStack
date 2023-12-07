<template>
    <MyCheckbox v-model="checked" @click="checkTask" :binary="true"></MyCheckbox>
</template>
  
  <script>
  
  import axios from 'axios';
  export default {
    name: "CheckBoxTask",
    props: {
        task: null,
    },
    emits: ["updatedTask"],
    data() {
      return {
        checked: null,
      };
    },
    mounted() {
        this.checked = Boolean(this.task.etat);
    },
    methods: {
        checkTask() {
            this.checked = !this.checked;
            const item = {
                    "nom": this.task.nom,
                    "description": this.task.description,
                    "priorite": this.task.priorite,
                    "etat": Number(this.checked),
                    "deadline": this.task.deadline,
                }
                const url = 'http://localhost:8080/api/v1/task/' + this.task.idTache;
                axios.put(url, item, {
                  headers: {
                      'Authorization': `Bearer ${sessionStorage.getItem('jwtToken')}`
                    }
                })
                .then(()=> {
                    this.$emit("updatedTask");
                })
                .catch((error) => {
                    console.log(error);
                });
        }
    },
  };
  </script>
  