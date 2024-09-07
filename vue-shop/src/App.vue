<template>
  <div id="app">
    <ListOfStores />
    <Statistics storeName="Store1" />
    <AuthForm />
  </div>
</template>

<script>
import ListOfStores from './components/ListOfStores.vue';
import Statistics from './components/Statistics.vue';
import AuthForm from './components/AuthForm.vue';

export default {
  name: 'App',
  components: {
    ListOfStores,
    Statistics,
    AuthForm,
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