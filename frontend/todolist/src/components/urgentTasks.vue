<template>
    <div class="mt-4">
        <h1 class="text-color"> Tâches urgentes restantes </h1>
        <TasksList :tasks="tasks" @updatedTask="initUrgentTasks" @delete-task="initUrgentTasks"></TasksList>
    </div>
    
</template>

<script>
import axios from 'axios';
import TasksList from './TasksList.vue';

export default {
    name: "UrgentTasks",
    emits: ["updateTask"],
    data() {
        return {
            nom: "Taches urgentes",
            tasks: [],
        };
    },
    mounted() {
        this.initUrgentTasks();
    },
    methods: {
        initUrgentTasks() {
            const token = sessionStorage.getItem('jwtToken');
            axios.get('http://localhost:8080/api/v1/tasks/urgent', {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            })
                .then((response) => {
                this.tasks = response.data.data;
                this.$emit("updateTask");
            })
                .catch((error) => {
                console.log(error);
            });
        }
    },
    components: { TasksList }
};
</script>
<style scoped></style>
