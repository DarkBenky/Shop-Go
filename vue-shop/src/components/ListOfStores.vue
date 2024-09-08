<template>
    <div>
        <nav>
            <div v-if="!show">
                <SearchBar :stores="stores" :items="currentStore" @selectStore="setCurrentStore" />
                <StoreMenu @selectStore="handleStoreSelection" />
                <!-- <ul>
                    <li v-for="store in stores" :key="store.id">
                        <a @click="setCurrentStore(store.name)" class="text-white">{{ store.name }}</a>
                    </li>
                </ul> -->
            </div>
            <transition name="fade">
                <div v-if="show">
                    <div class="hero">
                        <h1>{{ currentStoreName }}</h1>
                        <button @click="show = null">Close Store</button>
                    </div>
                    <div v-if="filteredStores != null">
                        <input type="text" v-model="query" @input="setQuery(query)" placeholder="Search for an item..."
                            class="search-input" />
                        <div>
                            <h1>Max Price</h1>
                            <input type="number" v-model="maxPrice" @change="setMaxPrice(maxPrice)"
                                class="search-input" />
                        </div>
                        <div>
                            <h1>Min Price</h1>
                            <input type="number" v-model="minPrice" @change="setMinPrice(minPrice)"
                                class="search-input" />
                        </div>
                        <div class="checkbox-container">
                            <label class="switch">
                                <input class="checkbox" type="checkbox" v-model="discount"
                                    @change="setDiscount(discount)" />
                                <span class="slider"></span>
                            </label>
                            <span class="checkbox-label">Discounted</span>
                        </div>
                        <div>
                            <h1>Category</h1>
                            <div class="category-container">
                                <div v-for="category in uniqueCategoriesWithImages" :key="category.name"
                                    class="category-item">
                                    <div v-if="!currentFilters.category.includes(category.name)"
                                        @click="setCategory(category.name)" class="category-unselected">
                                        <img :src="category.image" :alt="category.name" />
                                        <p>{{ category.name }}</p>
                                    </div>
                                    <div v-else class="category-selected">
                                        <img :src="category.image" :alt="category.name"
                                            @click="clearCategory(category.name)" />
                                        <p>{{ category.name }}</p>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <!-- Add Sorting Dropdown -->
                        <div class="sort-container">
                            <label for="sortOptions">Sort by:</label>
                            <select v-model="sortOption" @change="sortItems" id="sortOptions">
                                <option value="name-asc">Name (A-Z)</option>
                                <option value="name-desc">Name (Z-A)</option>
                                <option value="price-asc">Price (Low to High)</option>
                                <option value="price-desc">Price (High to Low)</option>
                            </select>
                        </div>


                        <button @click="filterStores()">Apply Filters</button>

                        <!-- Display selected filters -->
                        <div v-if="currentFilters.category.length > 0 ||
                            currentFilters.minPrice !== null ||
                            currentFilters.maxPrice !== null ||
                            currentFilters.query !== '' ||
                            discount !== false
                        " class="selected-filters">
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
                            <span v-if="discount">
                                Discounted: {{ discount }}
                                <button @click="setDiscount(false)">X</button>
                            </span>
                        </div>

                        <div class="scrollable-container">
                            <div v-for="item in filteredStores" :key="item.id">
                                <ItemCard class="item" :item="item" :open="item.id == openItemId"
                                    @open="openItem(item.id)" @close="closeItem" />
                            </div>
                        </div>
                    </div>
                    <div v-else>
                        <div class="no-items-message">
                            <h1>There are no items in this store.</h1>
                        </div>
                    </div>

                </div>
            </transition>
        </nav>
    </div>
</template>

<script>
import axios from 'axios';
import SearchBar from './SearchBar.vue';
import ItemCard from './ItemCard.vue';
import StoreMenu from './StoreMenu.vue';
import Cookies from 'js-cookie';

export default {
    components: {
        SearchBar,
        ItemCard,
        StoreMenu,
    },
    data() {
        return {
            currentStoreName: '',
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
                discount: false,
            },
            openItemId: null,
            discount: false, // Initialize this to sync with currentFilters
            sortOption: 'name-asc', // Default value
            userId: null,
        };
    },
    mounted() {
        // Check for the existence of userId in cookies
    this.userId = Cookies.get('userId');
    if (this.userId) {
        console.log('User ID from cookies:', this.userId);
    } else {
        console.error('User ID is not available in cookies.');
    }
    },
    computed: {
        uniqueCategoriesWithImages() {
            if (this.currentStore.length > 0 && this.currentStore !== null) {
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
    watch: {
        // Watch for changes in currentFilters.discount and update the discount data property
        'currentFilters.discount'(newVal) {
            this.discount = newVal;
        },
        // Watch for changes in discount and update currentFilters.discount
        discount(newVal) {
            this.currentFilters.discount = newVal;
        },
    },
    created() {
        this.fetchStores();
    },
    methods: {
        handleStoreSelection(store) {
            this.setCurrentStore(store);
            this.show = true;
        },
        async recordClick(itemId) {
            if (this.userId === undefined || this.userId === null) {
                console.error('User ID is not set. Cannot record click.');
                return;
            }
            else {
                console.log('Recording click for item:', itemId);
                console.log('User ID:', this.userId);
                try {
                    await axios.post('http://localhost:8080/click', null, {
                        params: {
                            store_name: this.currentStoreName, // Include store_name as a query parameter
                            user_id: this.userId,
                            item_id: itemId,
                            
                        },
                    });
                } catch (error) {
                    console.error('Error recording click:', error.response?.data || error.message); // More detailed error message
                }
            }
        },
        async recordSearch(query) {
            const url = `http://localhost:8080/searchItems`;
            try {
                console.log('Recording search for:', query);
                console.log('User ID:', this.userId);
                await axios.post(url, {
                    user_id: this.userId,
                    store_name: this.currentStoreName,
                    query,
                });
                console.log('Recorded search for:', query);
            } catch (error) {
                console.error('Error recording search:', error);
            }
        },
        setDiscount(discount) {
            this.discount = discount;
        },
        openItem(itemId) {
            this.openItemId = itemId;
            this.recordClick(itemId); // Record the click when an item is opened
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
        },
        setMinPrice(minPrice) {
            this.currentFilters.minPrice = minPrice;
            if (minPrice === null) this.minPrice = '';
        },
        setCategory(category) {
            this.currentFilters.category.push(category);
        },
        clearCategory(category) {
            this.currentFilters.category = this.currentFilters.category.filter(name => name !== category);
        },
        filterStores() {
            if (this.currentStore === null) return;
            this.filteredStores = this.currentStore.filter(item => {
                return (
                    (this.currentFilters.category.length === 0 || this.currentFilters.category.includes(item.category)) &&
                    (this.currentFilters.minPrice === null || item.price >= this.currentFilters.minPrice) &&
                    (this.currentFilters.maxPrice === null || item.price <= this.currentFilters.maxPrice) &&
                    (this.currentFilters.query === '' || item.name.toLowerCase().includes(this.currentFilters.query.toLowerCase())) &&
                    (this.currentFilters.discount === false || item.discount != 0)
                );
            });

            if (this.currentFilters.query !== '') {
                this.recordSearch(this.currentFilters.query);
            }
            this.sortItems(); // Apply sorting after filtering
        },
        sortItems() {
            if (!this.sortOption) return; // Exit if sortOption is not defined
            const [key, order] = this.sortOption.split('-');
            this.filteredStores.sort((a, b) => {
                if (key === 'name') {
                    return order === 'asc' ? a.name.localeCompare(b.name) : b.name.localeCompare(a.name);
                } else if (key === 'price') {
                    return order === 'asc' ? a.price - b.price : b.price - a.price;
                }
                return 0;
            });
        },
        setCurrentStore(store) {
            console.log('Current store:', store);
            this.fetchStoreData(store.name);
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
                const response = await axios.get(`http://localhost:8080/store/${store}`);
                this.currentStoreName = store;
                this.currentStore = response.data;
                this.filteredStores = this.currentStore;
                this.show = true;
                console.log('Current store:', this.currentStore);
                this.filterStores();
            } catch (error) {
                console.error('Error fetching store data:', error);
                this.filteredStores = [];
                this.currentStore = [];
            }
        },
    },
};
</script>

<style scoped>
/* Container for the sort dropdown */
.sort-container {
    display: flex;
    align-items: center;
    margin-top: 15px;
}

/* Label for the dropdown */
.sort-container label {
    font-size: 16px;
    font-weight: 600;
    margin-right: 10px;
}

/* Style for the dropdown select */
.sort-container select {
    background-color: #333;
    border: 1px solid #555;
    border-radius: 5px;
    padding: 10px;
    font-size: 16px;
    font-family: Arial, sans-serif;
    color: #e0e0e0;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
    transition: border-color 0.3s ease, box-shadow 0.3s ease;
}

/* Dropdown arrow */
.sort-container select::-ms-expand {
    display: none;
}

/* Focus state for the select */
.sort-container select:focus {
    border-color: #4caf50;
    box-shadow: 0 0 4px rgba(0, 123, 255, 0.2);
    outline: none;
}

/* Option styling */
.sort-container option {
    background-color: #333;
    color: #e0e0e0;
}

/* Custom arrow styling (for browsers that support it) */
.sort-container::after {
    content: '▼';
    font-size: 12px;
    color: #888;
    margin-left: 8px;
    pointer-events: none;
}

/* Transition Classes */
.fade-enter-active,
.fade-leave-active {
    transition: opacity 0.5s ease-in-out, transform 0.5s ease-in-out;
}

.fade-enter,
.fade-leave-to

/* .fade-leave-active in <2.1.8 */
    {
    opacity: 0;
    transform: translateY(20px);
}

/* General Styles */
body {
    background-color: #12121202;
    /* Dark background for the whole page */
    color: #e0e0e0;
    /* Light text color for contrast */
}

/* Container for the checkbox and label */
.checkbox-container {
    display: flex;
    align-items: center;
    margin-top: 10px;
}

/* Styles for the switch */
.switch {
    position: relative;
    display: inline-block;
    width: 60px;
    height: 34px;
}

.switch input {
    opacity: 0;
    width: 0;
    height: 0;
}

/* Slider styles */
.slider {
    position: absolute;
    cursor: pointer;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: #ccc;
    transition: .4s;
    border-radius: 34px;
}

/* Slider before the toggle is checked */
.slider:before {
    position: absolute;
    content: "";
    height: 26px;
    width: 26px;
    border-radius: 50%;
    left: 4px;
    bottom: 4px;
    background-color: white;
    transition: .4s;
}

/* When the toggle is checked */
input:checked+.slider {
    background-color: #4caf50;
}

/* Slider move when checked */
input:checked+.slider:before {
    transform: translateX(26px);
}

/* Optional: Add a nice shadow effect on the slider */
.slider {
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

/* Label styling */
.checkbox-label {
    font-size: 1rem;
    color: #e0e0e0;
    margin-left: 10px;
}

/* Hero Section (Store Name) */
.hero {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px;
    background: linear-gradient(135deg, #1e1e1e, #333);
    /* Gradient background for the hero section */
    border-radius: 10px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.5);
    color: #fff;
    margin-bottom: 16px;
}

.hero h1 {
    font-size: 24px;
    font-weight: 600;
    margin: 0;
    text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.3);
    /* Subtle shadow to enhance text readability */
}

.hero button {
    background-color: #e57373;
    color: #fff;
    border: none;
    padding: 8px 16px;
    border-radius: 5px;
    cursor: pointer;
    font-size: 14px;
    transition: background-color 0.3s ease, transform 0.3s ease;
}

.hero button:hover {
    background-color: #f44336;
    transform: scale(1.05);
    /* Slight zoom effect on hover */
}

/* Scrollable Container */
.scrollable-container {
    margin-top: 10px;
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
    gap: 0px;
    max-width: 85vw;
    max-height: 80vh;
    overflow-y: auto;
    overflow-x: hidden;
    padding: 10px;
    border-radius: 10px;
    background: #1e1e1e;
    /* Dark background for the container */
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.3);
    /* Modern shadow for depth */
}

.scrollable-container::-webkit-scrollbar {
    width: 12px;
    background-color: #1e1e1e;
}

.scrollable-container::-webkit-scrollbar-thumb {
    background: #44444433;
    border-radius: 10px;
    transition: background 0.3s ease;
}

.scrollable-container::-webkit-scrollbar-thumb:hover {
    background: #666;
    /* Slightly lighter thumb on hover */
}

.scrollable-container .item {
    background: #2c2c2c3f;
    /* Slightly lighter background for items */
    border-radius: 10px;
    padding: 10px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
    /* Light shadow for depth */
    transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.scrollable-container .item:hover {
    transform: translateY(-5px);
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.4);
    /* Elevation effect on hover */
}

.scrollable-container .item p {
    color: #e0e0e0;
    font-size: 14px;
    margin: 5px 0;
}

.scrollable-container .item a {
    display: inline-block;
    margin-top: 10px;
    padding: 8px 12px;
    color: #fff;
    background-color: #007bff;
    text-decoration: none;
    border-radius: 5px;
    transition: background-color 0.3s ease, transform 0.3s ease;
}

.scrollable-container .item a:hover {
    background-color: #0056b3;
    transform: scale(1.05);
    /* Slight zoom effect on hover */
}

/* Search Input */
.search-input {
    width: 100%;
    padding: 10px;
    margin-bottom: 10px;
    box-sizing: border-box;
    background-color: #333;
    /* Dark background for input */
    border: 1px solid #555;
    /* Slightly lighter border */
    color: #e0e0e0;
    /* Light text color */
    border-radius: 5px;
    /* Rounded corners */
}

.search-input::placeholder {
    color: #888;
    /* Lighter placeholder text */
}

/* Category Container */
.category-container {
    display: flex;
    overflow-x: auto;
    padding: 10px;
    border-radius: 12px;
    background: linear-gradient(135deg, #2c2c2c, #1f1f1f);
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.3);
    /* Dark gradient background with shadow for a modern look */
}

/* Category Items */
.category-item {
    margin-block: 10px;
    flex: 0 0 auto;
    margin-right: 15px;
    text-align: center;
    cursor: pointer;
    width: 130px;
    transition: transform 0.3s ease, box-shadow 0.3s ease;
    /* Added transitions for smooth scaling and shadow effects */
}

/* Category Item Image */
.category-item img {
    width: 100%;
    height: 130px;
    object-fit: cover;
    border-radius: 12px;
    margin-bottom: 5px;
    transition: transform 0.3s ease;
    /* Smooth image scaling */
}

/* Unselected Category Items */
.category-unselected {
    border-radius: 12px;
    opacity: 1;
    border: 2px solid transparent;
    background-color: #333;
    /* Dark background for unselected items */
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
    /* Slight shadow for unselected items */
}

.category-unselected:hover {
    border-color: #666;
    transform: scale(1.05);
    /* Slightly enlarges item on hover */
}

/* Selected Category Items */
.category-selected {
    border-radius: 12px;
    opacity: 1;
    border: 2px solid #00bcd4;
    background-color: #2c2c2c;
    /* Slightly lighter background for selected items */
    box-shadow: 0 6px 12px rgba(0, 0, 0, 0.4);
    /* More pronounced shadow for selected items */
    transform: scale(1.05);
    /* Slightly enlarges selected item */
}

/* Image Scaling Effect */
.category-item img:hover {
    transform: scale(1.05);
    /* Slightly enlarges image on hover */
}

/* Selected Filters */
.selected-filters {
    padding: 10px;
    background-color: #1f1f1f;
    /* Dark background for selected filters */
    border-radius: 10px;
}

.selected-filters span {
    display: inline-block;
    margin-right: 10px;
}

.selected-filters button {
    background-color: #e57373;
    /* Light red color for remove buttons */
    color: #fff;
    /* White text color */
    border: none;
    padding: 2px 8px;
    border-radius: 5px;
    cursor: pointer;
    font-size: 14px;
    margin-left: 5px;
}

.selected-filters button:hover {
    background-color: #f44336;
    /* Darker red on hover */
}

/* General Button */
button {
    background-color: #333;
    /* Dark background for buttons */
    color: #e0e0e0;
    /* Light text color */
    border: none;
    padding: 10px 20px;
    border-radius: 5px;
    cursor: pointer;
    font-size: 16px;
    transition: background-color 0.3s ease, transform 0.3s ease;
    margin-top: 10px;
}

button:hover {
    background-color: #555;
    transform: scale(1.05);
    /* Slightly lighter background on hover */
}

/* Suggestions */
.suggestions {
    list-style: none;
    padding: 0;
    margin: 0;
    border: 1px solid #333;
    /* Dark border for suggestions */
    background-color: #1f1f1f;
    /* Dark background for suggestions */
    max-height: 150px;
    overflow-y: auto;
}

.suggestions li {
    padding: 10px;
    cursor: pointer;
}

.suggestions li:hover {
    background-color: #333;
    /* Slightly lighter background on hover */
}

/* No Items Message */
.no-items-message {
    text-align: center;
    padding: 40px;
    color: #e57373;
    /* Light red color for the message text */
    font-size: 18px;
    font-weight: 600;
    background-color: #1f1f1f;
    /* Dark background to match the theme */
    border-radius: 10px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.5);
    /* Shadow for depth */
    margin-top: 20px;
    opacity: 0.9;
    transition: transform 0.3s ease, opacity 0.3s ease;
}

.no-items-message:hover {
    transform: scale(1.05);
    /* Slightly enlarges message on hover */
    opacity: 1;
    /* Increase opacity on hover */
}

.no-items-message h1 {
    margin: 0;
    padding: 0;
    color: #ffffff;
    /* Ensure the text is clearly visible */
    text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.3);
    /* Add a subtle text shadow for readability */
}
</style>
