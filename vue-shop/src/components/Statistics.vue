<template>
    <!-- Bar Chart Made of Divs -->
    <div v-if="groupedData">
        <div v-for="(items, key) in groupedData" :key="key" class="bar-container">
            <div class="text">
                <h2>{{ key }}</h2>
                <p>Number of Opened Items : {{ items.length }}</p>
            </div>
            <div v-for="item in items" :key="item.id" class="bar bar-loaded"
                :style="{ '--bar-height': getBarHeight(item) + '%' }">
                ***
            </div>
        </div>
    </div>
</template>

<script>
import axios from 'axios';

export default {
    data() {
        return {
            selectedPeriod: 'day',
            rawData: [],
            groupedData: {},
            barsLoaded: false // New data property to control animation
        };
    },
    props: {
        storeName: {
            type: String,
            required: true
        }
    },
    async mounted() {
        await this.fetchData();
        this.$nextTick(() => {
            this.barsLoaded = true; // Trigger animation after DOM update
        });
    },
    methods: {
        async fetchData() {
            const url = 'http://localhost:8080/statistics/';
            try {
                const response = await axios.get(url + this.storeName);
                this.rawData = response.data;
                console.log('Fetched data:', this.rawData);
                this.groupedData = this.groupByPeriod(this.selectedPeriod);
            } catch (error) {
                console.error('Error fetching data:', error);
            }
        },
        groupByPeriod(period) {
            return this.rawData.reduce((acc, item) => {
                const date = new Date(item.time_of_click);
                const key = date.toISOString().slice(0, period === 'day' ? 10 : 7);
                if (!acc[key]) {
                    acc[key] = [];
                }
                acc[key].push(item);
                return acc;
            }, {});
        },
        getBarHeight(item) {
            const maxValue = Math.max(...this.rawData.map(data => data.value));
            return (item.value / maxValue) * 100;
        }
    }
};
</script>

<style scoped>
.text {
    flex: 1;
}

.bar-container {
    display: flex;
    margin-top: 10px;
    margin-bottom: 10px;
    border-radius: 8px;
    overflow: hidden;
    padding: 10px;
    background: linear-gradient(135deg, #1e1e1e, #333);
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.bar-container h2 {
    font-size: 1.2em;
    margin: 0 0 10px 0;
    padding: 0;
    color: #e0e0e0;
}

.bar {
    background: linear-gradient(to top, hsl(122, 39%, 49%), #8bc34a);
    /* Gradient effect */
    color: rgba(255, 255, 255, 0);
    text-align: center;
    margin: 2px;
    border-radius: 5px;
    width: 5%;
    display: inline-block;
    height: 0;
    /* Initial height is 0 */
    transition: wid 0.5s ease-in-out;
    /* Smooth transition */
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.bar-container .bar-loaded {
    height: var(--bar-height);
    /* Use CSS variable for dynamic height */
}

.bar-container .bar:hover {
    background: linear-gradient(to top, #45a049, #7cb342);
}
</style>
