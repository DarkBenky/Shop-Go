<template>
    <div v-if="open" class="item-card-modal" @click="$emit('close')">
        <div class="modal-content" @click.stop>
            <button @click="$emit('close')" class="close-button">X</button>
            <div v-for="store in stores" :key="store.id">
                <div @click="selectStore(store)">
                    <h2>{{ store.name }}</h2>
                    <img :src="store.image" alt="Store Image" class="carousel-image" />
                    <p>{{ store.address }}</p>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    props: {
        stores: {
            type: Array, // Change type to Array if it's an array of stores
            required: true,
        },
        open: {
            type: Boolean,
            default: false,
        }
    },
    methods: {
        selectStore(store) {
            this.$emit('selectStore', store); // Emit the selected store object
        },
    },
};
</script>

<style scoped>
/* Scrollable Container */
.scrollable-container {
    display: grid;
    /* grid-template-columns: repeat(auto-fill, minmax(250px, 1fr)); */
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


/* Fade transition */
.fade-enter-active, .fade-leave-active {
    transition: opacity 0.2s ease;
}

.fade-enter, .fade-leave-to {
    opacity: 0;
}

/* Image Carousel Container with Arrows */
.image-carousel {
    display: flex;
    align-items: center;
    justify-content: center;
    flex-direction: row;
    position: relative;
    /* Adjust size as needed */
    margin: 0 auto;
    text-align: center;
}

.carousel-image {
    object-fit: cover;
    width: 95%;
    aspect-ratio: 16/9;
    border-radius: 10px;
    margin-bottom: 15px;
}

/* Arrow styles */
.arrow {
    position: absolute;
    top: 50%;
    transform: translateY(-50%);
    background: rgba(0, 0, 0, 0.5);
    color: white;
    border: none;
    padding: 10px;
    cursor: pointer;
    z-index: 10;
    font-size: 1.2rem;
    border-radius: 50%;  /* Ensure the arrows are round */
    transition: background-color 0.3s ease, transform 0.3s ease; /* Add transform transition */
}

.left-arrow {
    left: -15px;
}

.right-arrow {
    right: -15px;
}

.arrow:hover {
    background: rgba(0, 0, 0, 0.7);
}

/* Image fade-in animation */
.image-carousel .carousel-image-enter-active,
.image-carousel .carousel-image-leave-active {
    opacity: 0;
    transition: opacity 0.5s ease;
}

/* Other styles as previously defined */
/* Discount Section Styles */
.discount {
    margin-top: 10px;
    padding: 15px;
    background: linear-gradient(135deg, #ffe57f, #ffd54f);
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
    display: flex;
    flex-direction: column;
    align-items: center;
}

.discount h2 {
    margin: 5px 0;
    font-size: 1.2em;
    font-weight: 500;
}

.discount h2:first-of-type {
    color: #d32f2f;
}

.discount h2:last-of-type {
    color: #388e3c;
    font-weight: bold;
}

/* Item Card Styles */
.item-card {
    background: linear-gradient(135deg, #1f1f1f, #333);
    border: 1px solid #444;
    border-radius: 10px;
    padding: 15px;
    margin-bottom: 20px;
    text-align: center;
    cursor: pointer;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.3);
    transition: transform 0.3s ease, box-shadow 0.3s ease;
    width: 250px;
    height: 350px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: space-between;
}

.item-card:hover {
    transform: translateY(-5px);
    box-shadow: 0 6px 10px rgba(0, 0, 0, 0.4);
}

.item-card img {
    width: 100%;
    height: 150px;
    object-fit: cover;
    border-radius: 8px;
    margin-bottom: 15px;
}

/* Modal Background Styles */
.item-card-modal {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.4);
    backdrop-filter: blur(10px);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
    animation: fadeIn 0.3s ease-out;
}

/* Modal Content Styles */
.modal-content {
    overflow-y: scroll;
    overflow-x: hidden;
    background: rgba(44, 44, 44, 0.9);
    padding: 40px;
    border-radius: 20px;
    text-align: center;
    position: relative;
    height: 85%;
    width: 75%;
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.6);
    transform: scale(0.9);
    animation: zoomIn 0.3s ease-out forwards;
}

/* scroll bar */
.modal-content::-webkit-scrollbar {
    width: 10px;
}

.modal-content::-webkit-scrollbar-thumb {
    background: #888;
    border-radius: 5px;
}

.modal-content::-webkit-scrollbar-thumb:hover {
    background: #555;
}

.modal-content img {
    object-fit: cover;
    width: 95%;
    aspect-ratio: 16/9;
    border-radius: 10px;
    margin-bottom: 15px;
}

.modal-content h1,
.modal-content h2,
.modal-content p {
    color: #fff;
}

/* Close Button Styles */
.close-button {
    position: absolute;
    top: 10px;
    right: 10px;
    background: #d32f2f;
    color: #fff;
    border: none;
    font-size: 20px;
    cursor: pointer;
    padding: 8px 12px;
    border-radius: 50%;
    transition: background-color 0.3s ease, transform 0.3s ease;
}

.close-button:hover {
    background: #b71c1c;
    transform: scale(1.1);
}

/* Animations */
@keyframes fadeIn {
    from {
        opacity: 0;
    }

    to {
        opacity: 1;
    }
}

@keyframes zoomIn {
    from {
        transform: scale(0.9);
    }

    to {
        transform: scale(1);
    }
}

@keyframes zoomOut {
    from {
        transform: scale(1);
    }

    to {
        transform: scale(0.9);
    }
}
</style>
