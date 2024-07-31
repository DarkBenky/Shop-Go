<template>
  <div id="app">
    <ListOfStores />
  </div>
</template>

<script>
import ListOfStores from './components/ListOfStores.vue';

export default {
  name: 'App',
  components: {
     ListOfStores
  }, 
  data() {
    return {
      currentStore: null,
    };
  },
};
</script>

<style>
/* Add your custom styles here */
</style>