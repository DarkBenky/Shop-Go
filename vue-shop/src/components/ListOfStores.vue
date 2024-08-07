<template>
    <div>
        <nav>
            <div>
                <SearchBar :stores="stores" :items="currentStore" @selectStore="setCurrentStore" />
                <!-- <ul>
                    <li v-for="store in stores" :key="store.id">
                        <a @click="setCurrentStore(store.name)" class="text-white">{{ store.name }}</a>
                    </li>
                </ul> -->
            </div>

            <div v-if="show">
                <input type="text" v-model="query" @input="setQuery(query)" placeholder="Search for an item..."
                    class="search-input" />
                <div>
                    <h1>Max Price</h1>
                    <input type="number" v-model="maxPrice" @change="setMaxPrice(maxPrice)" class="search-input" />
                </div>
                <div>
                    <h1>Min Price</h1>
                    <input type="number" v-model="minPrice" @change="setMinPrice(minPrice)" class="search-input" />
                </div>
                <div>
                    <h1>Category</h1>
                    <div class="category-container">
                        <div v-for="category in uniqueCategoriesWithImages" :key="category.name" class="category-item">
                            <div v-if="!currentFilters.category.includes(category.name)"
                                @click="setCategory(category.name)">
                                <img :src="category.image" :alt="category.name" />
                                <p>{{ category.name }}</p>
                            </div>

                        </div>
                    </div>
                </div>

                <button @click="filterStores">Apply Filters</button>

                <!-- Display selected filters -->
                <div class="selected-filters">
                    <span v-if="currentFilters.category.length > 0">
                        <span v-for="name in currentFilters.category" :key="name">
                            Category: {{ name }}
                            <button @click="clearCategory(name)">X</button>
                        </span>

                    </span>
                    <span v-if="currentFilters.minPrice !== null">
                        Min Price: {{ minPrice }}
                        <button @click="setMinPrice(null)">X</button>
                    </span>
                    <span v-if="currentFilters.maxPrice !== null">
                        Max Price: {{ maxPrice }}
                        <button @click="setMaxPrice(null)">X</button>
                    </span>
                    <span v-if="currentFilters.query !== ''">
                        Query: {{ query }}
                        <button @click="setQuery('')">X</button>
                    </span>
                </div>

                <div>
                    <button @click="currentStore = null">X</button>
                    <div v-for="item in filteredStores" :key="item.id">
                        <ItemCard :item="item" :open="item.id === openItemId" @open="openItem(item.id)"
                            @close="closeItem" />
                    </div>
                </div>
            </div>
        </nav>
    </div>
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
            show: false,
            currentStore: [],
            filteredStores: [],
            stores: [],
            query: '',
            minPrice: null,
            maxPrice: null,
            currentFilters: {
                category: [],
                minPrice: null,
                maxPrice: null,
                query: '',
            },
            openItemId: null,
        };
    },
    computed: {
        uniqueCategoriesWithImages() {
            if (this.currentStore.length > 0) {
                const categories = this.currentStore.reduce((acc, item) => {
                    if (!acc[item.category]) {
                        acc[item.category] = item.image;
                    }
                    return acc;
                }, {});

                return Object.entries(categories).map(([name, image]) => ({
                    name,
                    image,
                }));
            }
            return [];
        },
    },
    created() {
        this.fetchStores();
    },
    methods: {
        openItem(itemId) {
            this.openItemId = itemId;
        },
        setMaxPrice(maxPrice) {
            this.currentFilters.maxPrice = maxPrice;
            if (maxPrice === null) this.maxPrice = '';
        },
        closeItem() {
            this.openItemId = null;
        },
        setQuery(query) {
            this.currentFilters.query = query;
            if (query === '') this.query = '';
            console.log('Query:', this.query);
        },
        setMinPrice(minPrice) {
            this.currentFilters.minPrice = minPrice;
            if (minPrice === null) this.minPrice = '';
            console.log('Min price:', this.minPrice);
        },
        setCategory(category) {
            this.currentFilters.category.push(category);
        },
        clearCategory(category) {
            this.currentFilters.category = this.currentFilters.category.filter(name => name !== category);
        },
        filterStores() {
            this.filteredStores = this.currentStore.filter(item => {
                return (
                    (this.currentFilters.category.length === 0 || this.currentFilters.category.includes(item.category)) &&
                    (this.currentFilters.minPrice === null || item.price >= this.currentFilters.minPrice) &&
                    (this.currentFilters.maxPrice === null || item.price <= this.currentFilters.maxPrice) &&
                    (this.currentFilters.query === '' || item.name.toLowerCase().includes(this.currentFilters.query.toLowerCase()))
                );
            });
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
                this.filteredStores = this.currentStore;
                this.show = true;
                console.log('Current store:', this.currentStore);
                this.filterStores();
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

.category-container {
    display: flex;
    overflow-x: auto;
    padding: 10px 0;
}

.category-item {
    flex: 0 0 auto;
    margin-right: 10px;
    text-align: center;
    cursor: pointer;
    width: 120px;
}

.category-item img {
    width: 100%;
    height: 120px;
    object-fit: cover;
    border-radius: 10px;
    margin-bottom: 5px;
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
