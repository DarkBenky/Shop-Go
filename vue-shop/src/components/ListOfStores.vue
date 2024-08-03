<template>
    <nav>
        <div>
            <SearchBar :stores="stores" :items="currentStore" @selectStore="setCurrentStore" />
            <ul>
                <li v-for="store in stores" :key="store.id">
                    <a @click="setCurrentStore(store.name)" class="text-white">{{ store.name }}</a>
                </li>
            </ul>
        </div>

        <div>
            <input type="text" v-model="query" @input="filterStores" placeholder="Search for a item..."
                class="search-input" />

            <div v-if="filteredStores.length == 0">
                <button @click="currentStore = null">X</button>
                <div v-for="item in currentStore" :key="item.id">
                    <ItemCard :item="item" />
                </div>
            </div>
            <div v-else>
                <button @click="currentStore = null">X</button>
                <div v-for="item in filteredStores" :key="item.id">
                    <h1>{{ item.name }}</h1>
                    <ItemCard :item="item" />
                </div>
            </div>
        </div>
    </nav>
</template>

<script>
import axios from 'axios';
import SearchBar from './SearchBar.vue';
import ItemCard from './ItemCard.vue';

export default {
    components: {
        SearchBar,
        ItemCard,
    },
    data() {
        return {
            currentStore: [],
            filteredStores: [],
            stores: [],
        };
    },
    created() {
        this.fetchStores();
    },
    methods: {
        filterStores() {
            if (this.currentStore.length > 0) {
                this.filteredStores = [];
                for (let i = 0; i < this.currentStore.length; i++) {
                    let item = this.currentStore[i];
                    for (let j = 0; j < Object.keys(item).length; j++) {
                        let value = String(Object.values(item)[j]);
                        if (value.toLowerCase().includes(this.query.toLowerCase())) {
                            this.filteredStores.push(item);
                            break;
                        }
                    }
                    console.log('Filtered stores:', this.filteredStores);
                }
            } if (this.query == '') {
                this.filteredStores = [];
            }
        },
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
.search-input {
    width: 100%;
    padding: 10px;
    margin-bottom: 10px;
    box-sizing: border-box;
}

.suggestions {
    list-style: none;
    padding: 0;
    margin: 0;
    border: 1px solid #ccc;
    max-height: 150px;
    overflow-y: auto;
}

.suggestions li {
    padding: 10px;
    cursor: pointer;
}

.suggestions li:hover {
    background-color: #f0f0f0;
}
</style>
