<template>
    <div v-if="tasks.length <= 0">
        <div class="text-center mb-2 boutonTabView">
            <MyButton icon="ri-add-fill" label="Ajouter une tache" severity="success" rounded @click="openAddTask(id)"></MyButton>
        </div>
        Aucunes tâches disponibles
    </div>
    
    <div v-else>
        <div class="text-right boutonTabView">
            <MyButton icon="ri-add-fill" label="Ajouter une tache" severity="success" rounded @click="openAddTask(id)"></MyButton>
        </div>
        <MyTabView>
            <MyTabPanel header="A Faire">
                <DataView :value="currentTasks">
                    <template #list="slotProps">
                        <div class="grid grid-nogutter">
                            <div v-for="(item, index) in slotProps.items" :key="index" class="col-12">
                                <div class="flex flex-column xl:flex-row xl:align-items-start p-1 mb-2 gap-4" :class="{ 'border-top-1 surface-border': index !== 0 }">
                                <div class="flex flex-column xl:flex-row gap-4 flex-1">
                                    <div class="flex align-items-center">
                                        <CheckBoxTask :task="item" @updatedTask="initTasks"></CheckBoxTask>
                                    </div>
                                    <div class="flex flex-column sm:flex-row justify-content-between align-items-center xl:align-items-start flex-1 gap-4">
                                    <div class="flex flex-column align-items-center sm:align-items-start gap-3">
                                        <div class="text-2xl font-bold text-900" v-tooltip="item.description">{{ item.nom }}</div>
                                        <div class="flex align-items-center gap-3">
                                        <span class="flex align-items-center gap-2">
                                            <i class="pi pi-calendar"></i>
                                            <span class="font-semibold">{{ item.deadline }}</span>
                                        </span>
                                        <MyTag :value="priorité[item.priorite]" :severity="getPriority(item)" />
                                        </div>
                                    </div>
                                    <div class="flex sm:flex-column align-items-center sm:align-items-end gap-3 sm:gap-2">
                                        <div class="datatable-buttons">
                                        <MyButton icon="ri-pencil-line" @click="openUpdateTask(item)" rounded />
                                        <MyButton icon="ri-delete-bin-line" @click="deleteTask(item)" rounded />
                                        </div>
                                    </div>
                                    </div>
                                </div>
                                </div>
                            </div>
                        </div>
                    </template>
                </DataView>
            </MyTabPanel>
            <MyTabPanel header="Toutes">
                <DataView :value="tasksToShow">
                    <template #list="slotProps">
                        <div class="grid grid-nogutter">
                            <div v-for="(item, index) in slotProps.items" :key="index" class="col-12">
                                <div class="flex flex-column xl:flex-row xl:align-items-start p-1 mb-2 gap-4" :class="{ 'border-top-1 surface-border': index !== 0, 'barre-ligne': item.etat === 1 }">
                                <div class="flex flex-column xl:flex-row gap-4 flex-1">
                                    <div class="flex align-items-center">
                                        <CheckBoxTask :task="item" @updatedTask="initTasks"></CheckBoxTask>
                                    </div>
                                    <div class="flex flex-column sm:flex-row justify-content-between align-items-center xl:align-items-start flex-1 gap-4">
                                    <div class="flex flex-column align-items-center sm:align-items-start gap-3">
                                        <div class="text-2xl font-bold text-900" :class="{ 'barre-nom': item.etat === 1 }" v-tooltip="item.description">{{ item.nom }}</div>
                                        <div class="flex align-items-center gap-3">
                                        <span class="flex align-items-center gap-2">
                                            <i class="pi pi-calendar"></i>
                                            <span class="font-semibold">{{ item.deadline }}</span>
                                        </span>
                                        <MyTag :value="priorité[item.priorite]" :severity="getPriority(item)" />
                                        </div>
                                    </div>
                                    <div class="flex sm:flex-column align-items-center sm:align-items-end gap-3 sm:gap-2">
                                        <div class="datatable-buttons">
                                        <MyButton icon="ri-pencil-line" @click="openUpdateTask(item)" rounded />
                                        <MyButton icon="ri-delete-bin-line" @click="deleteTask(item)" rounded />
                                        </div>
                                    </div>
                                    </div>
                                </div>
                                </div>
                            </div>
                        </div>
                    </template>
                </DataView>
            </MyTabPanel>
            <MyTabPanel header="Terminées">
                <DataView :value="finishedTask">
                    <template #list="slotProps">
                        <div class="grid grid-nogutter">
                            <div v-for="(item, index) in slotProps.items" :key="index" class="col-12">
                                <div class="flex flex-column xl:flex-row xl:align-items-start p-1 mb-2 gap-4" :class="{ 'border-top-1 surface-border': index !== 0 }">
                                <div class="flex flex-column xl:flex-row gap-4 flex-1">
                                    <div class="flex align-items-center">
                                        <CheckBoxTask :task="item" @updatedTask="initTasks"></CheckBoxTask>
                                    </div>
                                    <div class="flex flex-column sm:flex-row justify-content-between align-items-center xl:align-items-start flex-1 gap-4">
                                    <div class="flex flex-column align-items-center sm:align-items-start gap-3">
                                        <div class="text-2xl font-bold text-900" v-tooltip="item.description">{{ item.nom }}</div>
                                        <div class="flex align-items-center gap-3">
                                        <span class="flex align-items-center gap-2">
                                            <i class="pi pi-calendar"></i>
                                            <span class="font-semibold">{{ item.deadline }}</span>
                                        </span>
                                        <MyTag :value="priorité[item.priorite]" :severity="getPriority(item)" />
                                        </div>
                                    </div>
                                    <div class="flex sm:flex-column align-items-center sm:align-items-end gap-3 sm:gap-2">
                                        <div class="datatable-buttons">
                                        <MyButton icon="ri-pencil-line" @click="openUpdateTask(item)" rounded />
                                        <MyButton icon="ri-delete-bin-line" @click="deleteTask(item)" rounded />
                                        </div>
                                    </div>
                                    </div>
                                </div>
                                </div>
                            </div>
                        </div>
                    </template>
                </DataView>
            </MyTabPanel>
        </MyTabView>
    </div>
    <UpdateTaskPopup v-if="popup === 'UPDATETASK'" @close="closePopup" :task="taskToUpdate"/>
    <NewTaskPopup v-if="popup === 'NEWTASK'" @close="closePopup" :idTodolist="id"/>
</template>
  
<script>
    import axios from 'axios';
    import UpdateTaskPopup from './popup/UpdateTaskPopup.vue';
    import NewTaskPopup from './popup/NewTaskPopup.vue';
    import CheckBoxTask from './CheckBoxTask.vue';

    export default {
        name: 'TodolistTasks',
        components: { UpdateTaskPopup, NewTaskPopup, CheckBoxTask },
        props: {
            id: null,
        },
        data() {
            return {
                tasks: [],
                finishedTask: [],
                currentTasks: [],
                tasksToShow: [],
                priorité: ["Faible", "Moyenne", "Forte", "Urgente"],
                taskToUpdate: null,
                popup: null,
            };
        },
        mounted() {
            this.initTasks();
        },
        methods: {
            initTasks(){  
                const token = sessionStorage.getItem('jwtToken');
                const url = "http://localhost:8080/api/v1/todolist/" + this.id;
                axios.get(url, {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            })
                .then((response) => {
                    this.tasks = response.data.data;
                    this.finishedTask = this.tasks.filter(task => task.etat === 1);
                    this.currentTasks = this.tasks.filter(task => task.etat === 0);
                    this.tasksToShow = this.currentTasks.concat(this.finishedTask);
                })
                .catch((error) => {
                    console.log(error);
                });
            }, 
            deleteTask(task){
                const url = "http://localhost:8080/api/v1/todolist/" + this.id + "/task/" + task.idTache;
                axios.delete(url, {
                    headers: {
                        'Authorization': `Bearer ${sessionStorage.getItem('jwtToken')}`
                    }
                })
                .then(() => {
                    this.initTasks();
                })
                .catch((error) => {
                    console.log(error);
                });
            },
            openUpdateTask(task) {
                this.taskToUpdate = task;
                this.popup = "UPDATETASK";
            },
            getPriority(item) {
                switch (item.priorite) {
                    case 0: 
                        return 'success';

                    case 1:
                        return 'info';

                    case 2:
                        return 'warning';

                    case 3:
                        return 'danger';

                    default:
                        return null;
                }
            },
            closePopup(reload) {
                this.popup = null;
                if (reload) {
                    this.initTasks();
                }
            },
            openAddTask() {
                this.popup = "NEWTASK";
            },
        }
    }
</script>
    <style scoped>
    .boutonTabView {
        margin-bottom: 0;
    }
    .barre-nom {
        text-decoration: line-through;
        color: black;
        opacity: 0.5;
    }

    .barre-ligne .barre-nom {
        pointer-events: auto;
    }
    .barre-ligne {
        opacity: 0.5;
    }

    .barre-ligne .datatable-buttons {
        pointer-events: none;
        opacity: 0.5;
    }
</style>