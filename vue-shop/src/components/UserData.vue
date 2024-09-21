<template>
    <div class="user-details">
        <div class="hero">
            <h1 v-if="user">Hello, {{ user.username }}</h1>
        </div>
        <div v-if="user" class="user-info">
            <div class="user-info-item">
                <p><strong>Username:</strong></p>
                <input v-model="usernameInput" placeholder="Enter username" />
                <button @click="updateUsername">{{ user.username ? 'Modify' : 'Add' }} Username</button>
            </div>
            <div class="user-info-item">
                <p><strong>Email:</strong></p>
                <input v-model="emailInput" placeholder="Enter email" />
                <button @click="updateEmail">{{ user.email ? 'Modify' : 'Add' }} Email</button>
            </div>
            <div class="user-info-item">
                <p><strong>Address:</strong></p>
                <input v-model="addressInput" placeholder="Enter address" />
                <button @click="updateAddress">{{ user.address ? 'Modify' : 'Add' }} Address</button>
            </div>
            <div class="user-info-item">
                <p>Password:</p>
                <button @click="changePassword">Change Password</button>
            </div>
        </div>
        <div v-else-if="loading" class="loading-message">
            <div class="loader"></div>
            <p>Loading user data...</p>
        </div>
        <div v-else class="no-user-message">
            <p>User not found.</p>
        </div>
    </div>
</template>

<script>
import Cookies from 'js-cookie';

export default {
    data() {
        return {
            user: null,
            loading: true,
            userId: 1, // Replace with the actual user ID you want to fetch
            usernameInput: '',
            emailInput: '',
            addressInput: '',
        };
    },
    created() {
        this.fetchUser();
    },
    methods: {
        async fetchUser() {
            console.log("Cookies data:", Cookies.get());
            this.userId = Cookies.get('userId');
            try {
                const response = await fetch(`http://localhost:8080/user/${this.userId}`);
                if (response.ok) {
                    this.user = await response.json();
                    // Initialize inputs with current user data
                    this.usernameInput = this.user.username || '';
                    this.emailInput = this.user.email || '';
                    this.addressInput = this.user.address || '';
                } else {
                    this.user = null;
                }
            } catch (error) {
                console.error("Error fetching user:", error);
                this.user = null;
            } finally {
                this.loading = false;
            }
        },
        async updateUsername() {
            await this.updateUser({ username: this.usernameInput });
        },
        async updateEmail() {
            await this.updateUser({ email: this.emailInput });
        },
        async updateAddress() {
            await this.updateUser({ address: this.addressInput });
        },
        async updateUser(updatedData) {
            try {
                const response = await fetch(`http://localhost:8080/user/${this.userId}`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(updatedData)
                });
                if (response.ok) {
                    this.user = { ...this.user, ...updatedData }; // Update local user data
                } else {
                    console.error("Failed to update user data.");
                }
            } catch (error) {
                console.error("Error updating user data:", error);
            }
        },
        changePassword() {
            console.log("Change Password");
        }
    }
};
</script>

<style scoped>
.user-details {
    max-width: 800px;
    margin: 0 auto;
    padding: 20px;
    border-radius: 8px;
    background: #1e1e1e; /* Dark background */
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.5);
    color: #e0e0e0; /* Light text color */
}

.user-details h1 {
    font-size: 28px;
    font-weight: 700;
    margin-bottom: 20px;
    text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.3); /* Enhance readability */
}

.user-info {
    margin-top: 20px;
}

.user-info-item {
    margin-bottom: 15px;
}

.user-info-item p {
    font-size: 16px;
    font-weight: 600;
}

.user-info-item input {
    width: 100%;
    padding: 10px;
    background-color: #333;
    border: 1px solid #555;
    border-radius: 5px;
    color: #e0e0e0;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
    transition: border-color 0.3s ease, box-shadow 0.3s ease;
    box-sizing: border-box;
}

.user-info-item input::placeholder {
    color: #888;
}

.user-info-item input:focus {
    border-color: #4caf50;
    box-shadow: 0 0 4px rgba(0, 123, 255, 0.2);
    outline: none;
}

.user-info-item button {
    width: 100%;
    padding: 10px;
    background-color: #e57373;
    color: #fff;
    border: none;
    border-radius: 5px;
    cursor: pointer;
    font-size: 16px;
    transition: background-color 0.3s ease, transform 0.3s ease;
    margin-top: 5px;
}

.user-info-item button:hover {
    background-color: #f44336;
    transform: scale(1.03);
}
</style>
