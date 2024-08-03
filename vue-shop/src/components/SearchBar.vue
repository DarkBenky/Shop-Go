<template>
    <div>
        <input
            type="text"
            v-model="query"
            @input="filterStores"
            placeholder="Search for a store..."
            class="search-input"
        />
        <ul v-if="filteredStores.length" class="suggestions">
            <li v-for="store in filteredStores" :key="store.id" @click="selectStore(store.name)">
                {{ store.name }}
            </li>
        </ul>
    </div>
</template>

<script>
export default {
    props: {
        stores: Array,
        items : Array,
    },
    data() {
        return {
            query: '',
            filteredStores: [],
        };
    },
    methods: {
        filterStores() {
            this.filteredStores = this.stores.filter(store =>
                store.name.toLowerCase().includes(this.query.toLowerCase())
            );
        },
        selectStore(storeName) {
            this.$emit('selectStore', storeName);
            this.query = ''; // clear the search query
            this.filteredStores = []; // clear the suggestions
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
