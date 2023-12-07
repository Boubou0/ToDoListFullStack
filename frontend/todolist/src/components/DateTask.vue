<template>
    <div>
        <h1 class="text-color" v-if="date">Tâches du {{ date }}</h1>
        <h1 v-else>Sélectionnez une date</h1>
        <div v-if="tasks.length > 0">
            <TasksList :tasks="tasks" @deleteTask="initTasks" @updatedTask="updateDate(date)"/>
        </div>
        <div class="text-color" v-else>
            Aucun tâche prévue pour ce jour.
        </div>
    </div>
</template>
  
<script>
    import axios from 'axios';
    import TasksList from '@/components/TasksList.vue'

    export default {
        name: 'DateTask',
        components: { TasksList },
        emits: ["changedDate"],
        data() {
        return {
            tasks: [],
        };
    },
        props: {
            date: null,
        },
        methods: {
            updateDate() {
                this.initTasks(this.date);
            },
            convertDate(inputDate) {
                const parts = inputDate.split('/');
                const formattedDate = `${parts[0]}${parts[1]}${parts[2]}`;
                return formattedDate;
            },
            initTasks(date) {  
                const token = sessionStorage.getItem('jwtToken');
                const formattedDate = this.convertDate(date);
                const url = 'http://localhost:8080/api/v1/tasks/' + formattedDate;
                axios.get(url, {
                    headers: {
                        'Authorization': `Bearer ${token}`
                    }
                }).then((response) => {
                    this.tasks = response.data.data; 
                })
                .catch((error) => {
                    console.log(error);
                });
            },
        },
        watch: {
            date(newDate) {
                if(newDate !== undefined){
                    this.initTasks(newDate);
                }
            },
        }
    }
</script>
<style scoped>
</style>