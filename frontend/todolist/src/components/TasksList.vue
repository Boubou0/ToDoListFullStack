<template>
    <div class="p-8">
        <DataView :value="tasks">
            <template #list="slotProps">
                <div class="grid grid-nogutter">
                    <div v-for="(item, index) in slotProps.items" :key="index" class="col-12">
                        <div class="flex flex-column xl:flex-row xl:align-items-start p-1 mb-2 gap-4" :class="{ 'border-top-1 surface-border': index !== 0, 'barre-ligne': item.etat === 1 }">
                            <div class="flex flex-column xl:flex-row gap-4 flex-1">
                                <div class="flex align-items-center">
                                    <CheckBoxTask :task="item" @updatedTask="this.$emit(`updatedTask`);"></CheckBoxTask>
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
                                <div class="flex flex-col sm:flex-row justify-end gap-3 sm:gap-2 items-center sm:items-end">
                                    <div class="datatable-buttons flex gap-3">
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
    </div>
    <UpdateTaskPopup v-if="popup === 'UPDATETASK'" @close="closePopup" :task="taskToUpdate"/>
</template>
  
  <script>
  import axios from 'axios';
  import UpdateTaskPopup from './popup/UpdateTaskPopup.vue';
  import CheckBoxTask from './CheckBoxTask.vue';
  export default {
    name: 'TasksList',
    emits: ["deleteTask", "updatedTask"],
    components: { UpdateTaskPopup, CheckBoxTask },
    props: {
        tasks: {
            type: Array,
            default: () => [],
        },
    },
    data() {
        return {
            priorité: ["Faible", "Moyenne", "Forte", "Urgente"],
            taskToUpdate: null,
            popup: null,
        };
    },
    methods: {
        openUpdateTask(task) {
            this.taskToUpdate = task;
            this.popup = "UPDATETASK";
        },
        closePopup(reload) {
            this.popup = null;
            if (reload) {
                this.$emit("updatedTask");
            }
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
        deleteTask(task){
            let date = task.deadline;
            const url = "http://localhost:8080/api/v1/task/" + task.idTache;
            axios.delete(url, { headers: {
                    'Authorization': `Bearer ${sessionStorage.getItem('jwtToken')}`
                }})
            .then(() => {
                this.$emit("deleteTask", date, {
                
            });
            })
            .catch((error) => {
                console.log(error);
            });
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