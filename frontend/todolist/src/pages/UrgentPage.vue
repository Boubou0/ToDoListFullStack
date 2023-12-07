<template>
     <div class="box">
        <urgentTasks></urgentTasks>
    </div>
</template>
  
<script>
import urgentTasks from '@/components/urgentTasks.vue';
  export default {
    name: 'UrgentPage',
    components: { urgentTasks },
  };
  </script>