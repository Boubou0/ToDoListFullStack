<template>
    <div class="box">
        <h1 class="text-color">Mes listes de tâches</h1>
        <MyButton
                icon="ri-add-line"
                @click="openNewTodolist"
                label="Créer une Todolist"
        />
        <div class="p-7">
            <DataView :value="todolists">
                <template #list="slotProps">
                    <div v-for="item in slotProps.items" :key="item.idTodoLists">
                        <MyPanel toggleable class="mb-4">
                            <template #header>
                                <span v-if="item === editedItem" class="p-input-icon-right">
                                    <InputText v-model="item.nom" v-on:keyup.enter="onValidation(item)"/>
                                    <MyButton icon="ri-check-line" @click.stop="onValidation(item)" />
                                    <MyButton icon="ri-close-circle-line" severity="danger" @click.stop="cancelModification(item)" />
                                    <div v-if="error" class="col-12 text-red-500 mt-2">
                                        {{ error }}
                                    </div>
                                </span>
                                <div v-else class="data-item-title flex items-center" @click="editTodolist(item)">
                                    <span>{{ item.nom }}</span>
                                    <i class="ri-pencil-line text-gray-500 ml-2 cursor-pointer" @click.stop="editTodolist(item)"></i>
                                </div>                       
                            </template>
                            <TodolistTasks :id=item.idTodoLists />
                            <template #icons>
                                <div>
                                    <MyButton icon="ri-delete-bin-line" severity="danger" rounded outlined @click="deleteTodolist(item)"></MyButton>
                                </div>
                            </template>
                        </MyPanel> 
                    </div>
                </template>
            </DataView>
        </div>
    </div>
    <NewTodolistPopup v-if="popup === 'NEWTODOLIST'" @close="closePopup" />
  </template>
  
  <script>
  import axios from 'axios';
  import NewTodolistPopup from './popup/NewTodolistPopup.vue';
  import TodolistTasks from './TodolistTasks.vue';

  export default {
    name: 'TodoLists',
    components: { NewTodolistPopup, TodolistTasks },
    data() {
        return {
            todolists: [],
            popup: null,
            editedItem: null,
            idTodolist: null,
            error: null,
            nameBeforeUpdate: null,
        };
    },
    mounted() {
        this.getTodoLists();
    },
    methods: {
        getTodoLists() {
            const token = sessionStorage.getItem('jwtToken');
            axios.get('http://localhost:8080/api/v1/todolist', {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            })
            .then((response) => {
                this.todolists = response.data.data;
            })
            .catch((error) => {
                console.log(error);
            });
        },
        openNewTodolist() {
            this.popup = "NEWTODOLIST";
        },
        editTodolist(item) {
            this.editedItem = item;
            this.nameBeforeUpdate = item.nom;
        },
        cancelModification(item) {
            item.nom = this.nameBeforeUpdate;
            this.editedItem = null;
            this.nameBeforeUpdate = null;
        },
        closePopup(reload) {
            this.popup = null;
            if (reload) {
                this.getTodoLists();
            }
        },
        deleteTodolist(item) {
            const url = 'http://localhost:8080/api/v1/todolist/' + item.idTodoLists;
            const token = sessionStorage.getItem('jwtToken');
            axios.delete(url, {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            })
            .then(()=> {
                this.getTodoLists();
            })
            .catch((error) => {
                console.log(error);
            });
        },
        onValidation(item) {
            if (item.nom.trim() === "") {
                this.error = "Le nom ne peut pas être vide.";
                return;
            }
            if (this.nameBeforeUpdate === item.nom){
                this.editedItem = null;
                this.nameBeforeUpdate = null;
                return;
            }
            this.error = null;
            const url = 'http://localhost:8080/api/v1/todolist/' + item.idTodoLists;
            axios.put(url, item, {
            headers: {
                'Authorization': `Bearer ${sessionStorage.getItem('jwtToken')}`
              }
          })
            .then(()=> {
                this.editedItem = null;
                this.nameBeforeUpdate = null;
            })
            .catch((error) => {
                console.log(error);
            });
        }
    }
  }
</script>
<style scoped>
  .data-item-content {
    flex: 1;
  }
  
  .data-item-title {
    font-size: 18px;
    font-weight: bold;
    cursor: pointer;
  }
</style>