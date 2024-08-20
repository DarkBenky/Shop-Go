<template>
  <div id="app">
    <ListOfStores />
    <!-- <StoreMenu /> -->
  </div>
</template>

<script>
import ListOfStores from './components/ListOfStores.vue';
// import StoreMenu from './components/StoreMenu.vue';

export default {
  name: 'App',
  components: {
    ListOfStores,
    // StoreMenu
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