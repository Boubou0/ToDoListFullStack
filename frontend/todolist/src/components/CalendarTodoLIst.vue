<template>
    <div>
        <h1 class="text-color">Mon calendrier</h1>
        <MyCalendar v-model="date" inline/>
  </div>
</template>
  
  <script>
  export default {
    name: 'CalendarTodolist',
    emits: ["newDate"],
    data() {
        return {
            date: null,
        };
    },
    mounted() {},
    methods: {
      formatDate(date) {
        const day = date.getDate().toString().padStart(2, '0');
        const month = (date.getMonth() + 1).toString().padStart(2, '0');
        const year = date.getFullYear().toString();
        return `${day}/${month}/${year}`;
      }
    },
    watch: {
      date(newDate) {
        this.$emit("newDate", this.formatDate(newDate));
      },
    }
  }
  </script>