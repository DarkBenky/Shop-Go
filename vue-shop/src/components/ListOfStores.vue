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
            <input type="text" v-model="query" @input="filterStores" placeholder="Search for an item..." class="search-input" />
            <div>
                <h1>Max Price</h1>
                <input type="number" v-model.number="maxPrice" @change="filterStores" class="search-input" /> 
            </div>
            <div>
                <h1>Min Price</h1>
                <input type="number" v-model.number="minPrice" @change="filterStores" class="search-input" /> 
            </div>
            <!-- Dropdown to select category -->
            <div>
                <h1>Category</h1>
                <select class="search-input" v-model="selectedCategory" @change="filterStores">
                    <option value="">All Categories</option>
                    <option v-for="category in uniqueCategories" :key="category" :value="category">{{ category }}</option>
                </select>
            </div>

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
            query: '',
            selectedCategory: '',
            maxPrice: null,
            minPrice: null,
        };
    },
    computed: {
        uniqueCategories() {
            if (this.currentStore.length > 0) {
                const categories = this.currentStore.map(item => item.category);
                return [...new Set(categories)]; // Remove duplicates
            }
            return [];
        },
    },
    created() {
        this.fetchStores();
    },
    methods: {
        filterStores() {
            if (this.currentStore.length > 0) {
                this.filteredStores = this.currentStore.filter(item => {
                    const matchesQuery = item.name.toLowerCase().includes(this.query.toLowerCase());
                    const matchesCategory = this.selectedCategory === '' || item.category === this.selectedCategory;
                    const matchesMinPrice = this.minPrice === null || item.price >= this.minPrice;
                    const matchesMaxPrice = this.maxPrice === null || item.price <= this.maxPrice;
                    
                    return matchesQuery && matchesCategory && matchesMinPrice && matchesMaxPrice;
                });
            } else if (this.query === '') {
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
                this.filterStores(); // Ensure stores are filtered when data is fetched
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
