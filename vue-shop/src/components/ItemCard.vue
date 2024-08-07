<template>
    <div v-if="open" class="item-card-modal">
        <div class="modal-content">
            <button @click="$emit('close')" class="close-button">X</button>
            <h1>{{ item.category }}</h1>
            <img :src="item.image" :alt="item.name" />
            <h1>{{ item.name }}</h1>
            <h2>{{ item.price }} $</h2>
            <p>{{ item.description }}</p>
        </div>
    </div>
    <div v-else class="item-card">
        <img :src="item.image" :alt="item.name" @click="$emit('open')" />
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
.item-card {
    border: 1px solid #ddd;
    padding: 20px;
    margin-bottom: 20px;
    text-align: center;
    cursor: pointer;
}

.item-card img {
    max-width: 100%;
    height: auto;
    margin-bottom: 10px;
}

/* Modal styles */
.item-card-modal {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.8);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
}

.modal-content {
    background: #474747;
    padding: 20px;
    border-radius: 8px;
    text-align: center;
    position: relative;
    max-height: 80vh;
    max-width: 80vw;
}
.modal-content img {
    object-fit: cover;
    width: 50%;
    height: auto;
    aspect-ratio: 1 / 1; /* This keeps the image square */
    max-width: 100%;
    max-height: 100%;
}


.close-button {
    position: absolute;
    top: 10px;
    right: 10px;
    background: red;
    color: #fff;
    border: none;
    font-size: 18px;
    cursor: pointer;
    padding: 5px 10px;
    border-radius: 50%;
}
</style>
