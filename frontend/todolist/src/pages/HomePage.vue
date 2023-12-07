<template>
  <div v-if="isLargeScreen">
    <MySplitter class="p-8">
      <MySplitterPanel :size=60 :minSize=60 class="flex align-items-center justify-content-center">
        <CalendarTodoLIst @newDate="newDateSelected"/>
      </MySplitterPanel>
      <MySplitterPanel :size=40 :minSize=40 class="flex align-items-center justify-content-center">
        <DateTask :date="date"/>
      </MySplitterPanel>
    </MySplitter>
  </div>
  <div v-else class="box">
    <CalendarTodoLIst @newDate="newDateSelected"/>
    <DateTask :date="date"/>
  </div>
</template>

<script>
import CalendarTodoLIst from '@/components/CalendarTodoLIst.vue';
import DateTask from '@/components/DateTask.vue';

export default {
  components: { CalendarTodoLIst, DateTask },
  data() {
        return {
            date: null,
            isLargeScreen: window.innerWidth > 980,
        };
    },
  mounted() {
    window.addEventListener('resize', this.handleResize);
  },
  beforeUnmount() {
    window.removeEventListener('resize', this.handleResize);
  },
  methods: {
    handleResize() {
      this.isLargeScreen = window.innerWidth > 980;
    },
    newDateSelected(newDate){
      this.date = newDate;
    }
  }
};
</script>