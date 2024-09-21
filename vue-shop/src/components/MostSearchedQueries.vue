<template>
    <div class="chart-container">
        <!-- Move the chart type selection to the top -->
        <div class="chart-type-select">
            <label for="chart-type-select">Select Chart Type:</label>
            <select v-model="selectedChartType" id="chart-type-select" @change="updateChart">
                <option value="bar">Bar Chart</option>
                <option value="pie">Pie Chart</option>
            </select>
        </div>
        <!-- Chart section -->
        <div class="chart-wrapper">
            <!-- Display Bar or Pie chart based on selection -->
            <div v-if="chartData.labels.length > 0">
                <component :is="selectedChartType === 'bar' ? 'Bar' : 'Pie'" :data="chartData"
                    :options="chartOptions" />
            </div>
        </div>
    </div>
</template>

<script>
import { defineComponent, ref, onMounted } from 'vue';
import axios from 'axios';
import { Bar, Pie } from 'vue-chartjs';
import { Chart as ChartJS, Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale, ArcElement } from 'chart.js';

// Register Chart.js components
ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale, ArcElement);

export default defineComponent({
    name: 'SearchedQueriesChart',
    components: {
        Bar,
        Pie
    },
    props: {
        storeName: {
            type: String,
            required: true
        }
    },
    setup(props) {
        const selectedChartType = ref('bar');
        const rawData = ref([]);
        const topSearches = ref([]);
        const chartData = ref({
            labels: [],
            datasets: [{
                label: 'Most Searched Queries',
                data: [],
                backgroundColor: ['#8bc34a', '#ff6384', '#36a2eb', '#ffcd56', '#4caf50', '#f44336'],
                borderColor: ['#f7e6a2'],
                borderWidth: 0
            }]
        });

        const chartOptions = {
            responsive: true,
            plugins: {
                legend: {
                    display: true,
                    position: 'top'
                },
                tooltip: {
                    callbacks: {
                        label: function (tooltipItem) {
                            return `Searches: ${tooltipItem.raw}`;
                        }
                    }
                }
            }
        };

        const fetchData = async () => {
            const url = 'http://localhost:8080/most-searched-items/';
            try {
                const response = await axios.get(url + props.storeName);
                rawData.value = response.data;
                console.log('Fetched data:', rawData.value);
                updateChartData();
            } catch (error) {
                console.error('Error fetching data:', error);
            }
        };

        const updateChartData = () => {
            const searchQueries = rawData.value.map(item => item.query);
            const searchCounts = rawData.value.map(item => item.count);

            chartData.value = {
                labels: searchQueries,
                datasets: [{
                    ...chartData.value.datasets[0],
                    data: searchCounts
                }]
            };

            // Populate topSearches for the list on the side
            topSearches.value = rawData.value.slice(0, 10); // Limit to top 10
        };

        onMounted(async () => {
            await fetchData();
        });

        return {
            selectedChartType,
            chartData,
            chartOptions,
            topSearches,
            updateChartData
        };
    }
});
</script>
<style scoped>
.chart-container {
    display: flex;
    flex-direction: column; /* Stack items vertically */
    background: linear-gradient(135deg, #2c2c2c, #1f1f1f);
    border-radius: 12px;
    padding: 20px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.3);
    color: #e0e0e0;
    margin-top: 20px;
    max-width: 800px;
}

.chart-type-select {
    display: flex;
    align-items: center;
    margin-bottom: 20px; /* Adjust spacing */
}

.chart-type-select label {
    font-size: 16px;
    font-weight: 600;
    margin-right: 10px;
}

.chart-type-select select {
    background-color: #333;
    border: 1px solid #555;
    border-radius: 5px;
    padding: 10px;
    font-size: 16px;
    color: #e0e0e0;
}

.chart-wrapper {
    width: 100%;
}
</style>