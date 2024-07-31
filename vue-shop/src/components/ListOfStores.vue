<!-- src/components/NavBar.vue -->
<template>
    <nav>
        <div>
            <ul>
                <li v-for="store in stores" :key="store.id">
                    <a @click="setCurrentStore(store.name)" class="text-white">{{ store.name }}</a>
                </li>
            </ul>
        </div>
        <div v-if="currentStore != null">
            <button @click="currentStore = null">X</button>
            <div v-for="item in currentStore" :key="item.id">
                <img :src="item.image" :alt="item.name" />
                <h1>{{ item.name }}</h1>
                <h2>{{ item.price }} $</h2>
                <p>{{ item.description }}</p>
            </div>
        </div>
    </nav>
</template>

<script>
import axios from 'axios';

export default {
    data() {
        return {
            currentStore: null,
            stores: [],
        };
    },
    created() {
        this.fetchStores();
    },
    methods: {
        setCurrentStore(store) {
            console.log('Current store:', store);
            this.fetchStoreData(store);
        },
        async fetchStores() {
            try {
                const response = await axios.get('http://localhost:8080/stores');
                this.stores = response.data;
            } catch (error) {
                console.error('Error fetching stores:', error);
            }
        },
        async fetchStoreData(store) {
            try {
                const response = await axios.get(`http://localhost:8080/${store}`);
                this.currentStore = response.data;
                console.log('Current store:', this.currentStore);
            } catch (error) {
                console.error('Error fetching store data:', error);
            }
        },
    },
};
</script>

<style scoped>
/* Add any styles you need here */
</style>