<template>
    <div v-if="open" class="item-card-modal" @click="$emit('close')">
        <div class="modal-content" @click.stop>
            <button @click="$emit('close')" class="close-button">X</button>
            <h1>{{ item.category }}</h1>
            <img :src="item.image" :alt="item.name" />
            <h1>{{ item.name }}</h1>
            <h2>{{ item.price }} $</h2>
            <div class="discount" v-if="item.discount != 0">
                <h2>Discount: {{ item.discount }} %</h2>
                <h2>Price with discount: {{ item.price - (item.price * item.discount) / 100 }} $</h2>
            </div>
            <p v-else>{{ item.description }}</p>
        </div>
    </div>
    <div v-else class="item-card" @click="$emit('open')">
        <img :src="item.image" :alt="item.name" />
        <h1>{{ item.name }}</h1>
        <div class="discount" v-if="item.discount != 0">
            <h2>Discount: {{ item.discount }} %</h2>
            <h2>Price with discount: {{ item.price - (item.price * item.discount) / 100 }} $</h2>
        </div>
        <h2 v-else>Price {{ item.price }} $</h2>
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
        }
    },
    data() {
        return {
            images: NaN,
        };
    },
    watch: {
        open(newValue) {
            if (newValue) {
                this.fetchImage();
            }
        }
    },
    methods: {
        async fetchImage() {
            const url = `http://localhost:8080/images/${this.item.id}`;
            console.log('Fetching image from URL:', url);

            try {
                const response = await fetch(url);
                if (response.ok) {
                    const data = await response.json();
                    this.images = data;
                    console.log('Fetched image:', this.images);
                } else {
                    console.error('Failed to fetch image:', response.statusText);
                }
            } catch (error) {
                console.error('Error fetching image:', error);
            }
        }
    },
};
</script>

<style scoped>
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
