<template>
  <div id="app">
    <RealTimeData />
    <WarningData />
  </div>
</template>

<script>
import RealTimeData from './components/RealTimeData.vue';
import WarningData from './components/WarningData.vue';

export default {
  name: 'App',
  components: {
    RealTimeData,
    WarningData,
  },
};
</script>