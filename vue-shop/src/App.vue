<template>
  <div id="app">
    <ListOfStores />
    <Statistics storeName="Store1" />
  </div>
</template>

<script>
import ListOfStores from './components/ListOfStores.vue';
import Statistics from './components/Statistics.vue';

export default {
  name: 'App',
  components: {
    ListOfStores,
    Statistics,
  },
  data() {
    return {
    };
  },
};
</script>

<style>
/* Add your custom styles here */
</style>