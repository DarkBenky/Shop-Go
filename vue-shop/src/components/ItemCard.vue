<template>
    <div v-if="open" class="item-card">
        <h1>{{ item.category }}</h1>
        <img :src="item.image" :alt="item.name" />
        <h1>{{ item.name }}</h1>
        <h2>{{ item.price }} $</h2>
        <p>{{ item.description }}</p>
    </div>
    <div v-else>
        <img :src="item.image" :alt="item.name" />
        <h1>{{ item.name }}</h1>
        <h2>{{ item.price }} $</h2>
    </div>
</template>

<script>
export default {
    props: {
        item: {
            type: Object,
            required: true,
        },
        open: {
            type: Boolean,
            default: false,
        },
    },
};
</script>

<style scoped>
.item-card-opened {
    border: 1px solid #ddd;
    padding: 20px;
    margin-bottom: 20px;
    text-align: center;
}

.item-card {
    border: 1px solid #ddd;
    padding: 20px;
    margin-bottom: 20px;
    text-align: center;
}

.item-card img {
    max-width: 100%;
    height: auto;
    margin-bottom: 10px;
}
</style>
