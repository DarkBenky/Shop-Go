<template>
  <div id="app">
    <ListOfStores />
    <Statistics storeName="Store1" />
    <AuthForm />
    <MostSearchedQueries storeName="Store1" />
  </div>
</template>

<script>
import ListOfStores from './components/ListOfStores.vue';
import Statistics from './components/Statistics.vue';
import AuthForm from './components/AuthForm.vue';
import MostSearchedQueries from './components/MostSearchedQueries.vue';

export default {
  name: 'App',
  components: {
    ListOfStores,
    Statistics,
    AuthForm,
    MostSearchedQueries,
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