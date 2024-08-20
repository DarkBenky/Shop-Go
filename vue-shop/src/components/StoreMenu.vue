<template>
    <div class="store-menu">
        <transition name="fade" @before-enter="beforeEnter" @enter="enter" @leave="leave">
            <div v-if="uniqueCategories.length > 0 && selectedCategory === ''" key="categories" class="category-container">
                <div
                    v-for="category in uniqueCategories"
                    :key="category.name"
                    class="category-item"
                    @click="filterStoresByCategory(category.name)"
                >
                    <img :src="category.image" :alt="category.name" />
                    <p>{{ category.name }}</p>
                </div>
            </div>
        </transition>

        <transition name="fade" @before-enter="beforeEnter" @enter="enter" @leave="leave">
            <ModalStores 
                v-if="selectedCategory !== ''" 
                @selectStore="handleStoreSelection" 
                :stores="filteredStores" 
                :open="true" 
                @close="clearSelection" 
            />
            <!-- TODO: Emit data to the List of Stores -->
        </transition>
    </div>
</template>

<script>
import axios from 'axios';
import ModalStores from './ModalStores.vue';

export default {
    components: {
        ModalStores,
    },
    data() {
        return {
            stores: [],
            selectedCategory: '',
            filteredStores: [],
        };
    },
    computed: {
        uniqueCategories() {
            const uniqueCategoriesMap = {};
            this.stores.forEach(store => {
                if (store.category && !uniqueCategoriesMap[store.category]) {
                    uniqueCategoriesMap[store.category] = {
                        name: store.category,
                        image: store.image,
                    };
                }
            });
            return Object.values(uniqueCategoriesMap);
        },
    },
    methods: {
        async fetchStores() {
            try {
                const response = await axios.get('http://localhost:8080/stores');
                this.stores = response.data;
            } catch (error) {
                console.error('Failed to fetch stores:', error);
            }
        },
        filterStoresByCategory(category) {
            this.selectedCategory = category;
            this.filterStores();
        },
        filterStores() {
            this.filteredStores = this.stores.filter(store => store.category === this.selectedCategory);
        },
        clearSelection() {
            this.selectedCategory = '';
        },
        handleStoreSelection(store) {
            // Handle the store selection, e.g., show details or perform an action
            console.log('Selected store:', store);
        },
        beforeEnter(el) {
            el.style.opacity = 0;
        },
        enter(el, done) {
            el.offsetHeight; // trigger reflow
            el.style.transition = 'opacity 0.5s ease';
            el.style.opacity = 1;
            done();
        },
        leave(el, done) {
            el.style.transition = 'opacity 0.5s ease';
            el.style.opacity = 0;
            done();
        }
    },
    mounted() {
        this.fetchStores();
    },
};
</script>


<style scoped>




/* Transition Classes */
.fade-enter-active, .fade-leave-active {
    transition: opacity 0.5s ease-in-out, transform 0.5s ease-in-out;
}

.fade-enter, .fade-leave-to /* .fade-leave-active in <2.1.8 */ {
    opacity: 0;
    transform: translateY(20px);
}

/* Styles for the main container */
.store-menu {
    padding: 20px;
    background: linear-gradient(135deg, #1e1e1e, #333);
    border-radius: 12px;
    box-shadow: 0 6px 12px rgba(0, 0, 0, 0.5);
    color: #e0e0e0;
    margin: 0 auto;
    max-width: 1200px;
}

/* Styles for the categories */
.category-container {
    display: flex;
    flex-wrap: wrap;
    gap: 15px;
    justify-content: center;
}

.category-item {
    width: 160px;
    cursor: pointer;
    transition: transform 0.3s ease, box-shadow 0.3s ease;
    border-radius: 12px;
    background-color: #2c2c2c;
    overflow: hidden;
}

.category-item img {
    width: 100%;
    height: 120px;
    object-fit: cover;
    border-radius: 12px 12px 0 0;
}

.category-item p {
    text-align: center;
    padding: 10px;
    color: #e0e0e0;
    font-size: 16px;
    background-color: #333;
}

.category-item:hover {
    transform: scale(1.05);
    box-shadow: 0 6px 12px rgba(0, 0, 0, 0.4);
}

/* Styles for the filtered stores list */
.store-list {
    margin-top: 20px;
}

.store-container {
    display: flex;
    flex-wrap: wrap;
    gap: 20px;
    justify-content: center;
}

.store-item {
    width: 220px;
    background: #2c2c2c;
    border-radius: 12px;
    box-shadow: 0 6px 12px rgba(0, 0, 0, 0.3);
    overflow: hidden;
    transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.store-item:hover {
    transform: scale(1.03);
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.4);
}

.store-image {
    width: 100%;
    object-fit: cover;
}

.store-info {
    padding: 15px;
}

.store-info h3 {
    font-size: 20px;
    margin: 0;
    color: #e0e0e0;
}

.store-info p {
    font-size: 14px;
    color: #b0b0b0;
    margin: 5px 0 0;
}

/* Styles for the back button */
.back-button {
    margin-top: 20px;
    padding: 12px 24px;
    background-color: #555;
    color: #e0e0e0;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    font-size: 16px;
    transition: background-color 0.3s ease, transform 0.3s ease;
}

.back-button:hover {
    background-color: #777;
    transform: scale(1.05);
}
</style>