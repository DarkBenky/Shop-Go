<template>
    <div>
        <input type="text" v-model="query" @input="filterStores" placeholder="Search for a store..."
            class="search-input" />
        <ul v-if="filteredStores.length && query != ''" class="suggestions">
            <li v-for="store in filteredStores" :key="store.id" @click="selectStore(store.name)">
                {{ store.name }}
            </li>
        </ul>
    </div>
</template>

<script>
import Cookies from 'js-cookie';
import axios from 'axios';
export default {
    props: {
        stores: Array,
        items: Array,
    },
    data() {
        return {
            query: '',
            filteredStores: [],
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
    methods: {
        filterStores() {
            this.filteredStores = this.stores.filter(store =>
                store.name.toLowerCase().includes(this.query.toLowerCase())
            );
        },
        async recordSearch(storeName) {
            const url = `http://localhost:8080/search`;
            try {
                console.log('Recording search for:', storeName);
                console.log('User ID:', this.userId);
                await axios.post(url, {
                    user_id: this.userId,
                    store_name: storeName,
                });
                console.log('Recorded search for:', storeName);
            } catch (error) {
                console.error('Error recording search:', error);
            }
        },
        selectStore(storeName) {
            this.recordSearch(storeName);
            this.$emit('selectStore', { name: storeName, id: this.stores.find(store => store.name === storeName).id });
            this.query = ''; // clear the search query
            this.filteredStores = []; // clear the suggestions
        },
    },
};
</script>

<style scoped>
.search-input {
    width: 100%;
    padding: 12px 15px;
    margin-bottom: 15px;
    box-sizing: border-box;
    background: #333;
    /* Dark background */
    border: 1px solid #555;
    /* Subtle border */
    border-radius: 8px;
    /* Rounded corners */
    color: #e0e0e0;
    /* Light text color */
    font-size: 16px;
    transition: background 0.3s ease, border-color 0.3s ease;
    /* Smooth transitions for background and border changes */
}

.search-input::placeholder {
    color: #888;
    /* Lighter placeholder text */
}

/* Suggestions Styles */
.suggestions {
    list-style: none;
    padding: 0;
    margin: 0;
    border: 1px solid #555;
    /* Slightly darker border */
    border-radius: 8px;
    /* Rounded corners */
    background: #1f1f1f;
    /* Dark background */
    max-height: 200px;
    /* Increased max-height */
    overflow-y: auto;
    /* Scroll if needed */
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.3);
    /* Subtle shadow */
    /* Modern shadow effect */
}

.suggestions li {
    padding: 12px;
    cursor: pointer;
    transition: background-color 0.3s ease, color 0.3s ease;
    /* Smooth transition for background color and text color */
}

.suggestions li:hover {
    background-color: #333;
    /* Slightly lighter background on hover */
    color: #e0e0e0;
    /* Light text color */
    /* Ensures good contrast */
}
</style>
