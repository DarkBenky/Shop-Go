<template>
    <div class="chart-container">
        <div class="chart-header">
            <div class="period-select">
                <label for="period-select">Select Period:</label>
                <select v-model="selectedPeriod" id="period-select" @change="updateChartData">
                    <option value="day">Day</option>
                    <option value="week">Week</option>
                    <option value="month">Month</option>
                </select>
            </div>
        </div>
        <div v-if="chartData.labels.length > 0" class="chart-wrapper">
            <Line :data="chartData" :options="chartOptions" />
        </div>
    </div>
</template>

<script>
import { defineComponent, ref, onMounted, watch } from 'vue';
import axios from 'axios';
import { Line } from 'vue-chartjs';
import { Chart as ChartJS, Title, Tooltip, Legend, LineElement, CategoryScale, LinearScale, PointElement } from 'chart.js';

// Register Chart.js components
ChartJS.register(Title, Tooltip, Legend, LineElement, CategoryScale, LinearScale, PointElement);

export default defineComponent({
    name: 'LineChartComponent',
    components: {
        Line
    },
    props: {
        storeName: {
            type: String,
            required: true
        }
    },
    setup(props) {
        const selectedPeriod = ref('day');
        const rawData = ref([]);
        const chartData = ref({
            labels: [],
            datasets: [{
                label: 'Store Views',
                data: [],
                borderColor: '#8bc34a',
                backgroundColor: 'rgba(139, 195, 74, 0.2)',
                borderWidth: 2,
                pointBackgroundColor: '#8bc34a',
                pointBorderColor: '#fff',
                pointBorderWidth: 1,
                pointRadius: 3
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
                            return `Views : ${tooltipItem.raw}`;
                        }
                    }
                }
            },
            scales: {
                x: {
                    title: {
                        display: true,
                        text: 'Time'
                    }
                },
                y: {
                    title: {
                        display: true,
                        text: 'Store Views'
                    },
                    beginAtZero: true
                }
            }
        };

        const fetchData = async () => {
            const url = 'http://localhost:8080/statistics/';
            try {
                const response = await axios.get(url + props.storeName);
                rawData.value = response.data;
                console.log('Fetched data:', rawData.value);
            } catch (error) {
                console.error('Error fetching data:', error);
            }
        };

        const getWeekNumber = (date) => {
            const d = new Date(Date.UTC(date.getFullYear(), date.getMonth(), date.getDate()));
            const dayNum = d.getUTCDay() || 7;
            d.setUTCDate(d.getUTCDate() + 4 - dayNum);
            const yearStart = new Date(Date.UTC(d.getUTCFullYear(), 0, 1));
            return Math.ceil((((d - yearStart) / 86400000) + 1) / 7);
        };

        const groupByPeriod = (period) => {
            console.log('Grouping by period:', period);
            const grouped = rawData.value.reduce((acc, item) => {
                const date = new Date(item.time_of_click);
                let key;
                let weekNum; // Move the declaration outside of the switch statement

                switch (period) {
                    case 'day':
                        key = date.toISOString().slice(0, 10); // YYYY-MM-DD
                        break;
                    case 'week':
                        weekNum = getWeekNumber(date); // Assign the value here
                        key = `${date.getFullYear()}-W${weekNum.toString().padStart(2, '0')}`;
                        break;
                    case 'month':
                        key = date.toISOString().slice(0, 7); // YYYY-MM
                        break;
                }

                if (!acc[key]) {
                    acc[key] = 0;
                }
                acc[key] += 1; // Count logs
                return acc;
            }, {});

            console.log('Grouped data:', grouped);
            return grouped;
        };

        const updateChartData = () => {
            const groupedData = groupByPeriod(selectedPeriod.value);
            chartData.value = {
                labels: Object.keys(groupedData),
                datasets: [{
                    ...chartData.value.datasets[0],
                    data: Object.values(groupedData)
                }]
            };
        };

        onMounted(async () => {
            await fetchData();
            updateChartData();
        });

        watch(selectedPeriod, () => {
            updateChartData();
        });

        return {
            selectedPeriod,
            chartData,
            chartOptions,
            updateChartData
        };
    }
});
</script>



<style scoped>
.chart-container {
    background: linear-gradient(135deg, #2c2c2c, #1f1f1f);
    border-radius: 12px;
    padding: 20px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.3);
    color: #e0e0e0;
    margin-top: 20px;
}

.chart-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
}

.chart-header h2 {
    font-size: 24px;
    font-weight: 600;
    margin: 0;
    text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.3);
}

.period-select {
    display: flex;
    align-items: center;
}

.period-select label {
    font-size: 16px;
    font-weight: 600;
    margin-right: 10px;
}

.period-select select {
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

.period-select select:focus {
    border-color: #4caf50;
    box-shadow: 0 0 4px rgba(76, 175, 80, 0.2);
    outline: none;
}

.period-select select option {
    background-color: #333;
    color: #e0e0e0;
}

.chart-wrapper {
    height: 400px;
    position: relative;
}
</style>