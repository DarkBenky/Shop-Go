<template>
    <div class="auth-form">
        <h2>{{ isLogin ? 'Login' : 'Sign Up' }}</h2>
        <form @submit.prevent="handleSubmit">
            <div>
                <label for="username-email">{{ isLogin ? 'Username/Email' : 'Username' }}</label>
                <input type="text" v-model="form.usernameOrEmail"
                    :placeholder="isLogin ? 'Username or Email' : 'Username'" required />
            </div>
            <div v-if="!isLogin">
                <label for="email">Email</label>
                <input type="email" v-model="form.email" placeholder="Email" required />
            </div>
            <div>
                <label for="password">Password</label>
                <input type="password" v-model="form.password" placeholder="Password" required />
            </div>
            <button type="submit">{{ isLogin ? 'Login' : 'Sign Up' }}</button>
        </form>
        <button @click="toggleAuthMode">
            Switch to {{ isLogin ? 'Sign Up' : 'Login' }}
        </button>
    </div>
</template>

<script>
import axios from 'axios';
import Cookies from 'js-cookie';

export default {
    data() {
        return {
            isLogin: true,
            form: {
                usernameOrEmail: '',
                email: '',
                password: ''
            },
            token: null
        };
    },
    methods: {
        handleSubmit() {
            const url = this.isLogin
                ? 'http://localhost:8080/login'
                : 'http://localhost:8080/register';

            const data = this.isLogin
                ? { usernameOrEmail: this.form.usernameOrEmail, password: this.form.password }
                : { username: this.form.usernameOrEmail, email: this.form.email, password: this.form.password };

            axios.post(url, data)
                .then(response => {
                    this.token = response.data.token;
                    Cookies.set('authToken', this.token, { expires: 7 });
                    console.log('Successfully logged in or registered');
                })
                .catch(error => {
                    console.error('Error during authentication', error);
                });
        },
        toggleAuthMode() {
            this.isLogin = !this.isLogin;
        },
        checkLoggedIn() {
            // Check if there's an auth token in cookies
            const token = Cookies.get('authToken');
            if (token) {
                this.token = token;
                console.log('User is already logged in');
            } else {
                console.log('No user is logged in');
            }
        }
    },
    mounted() {
        this.checkLoggedIn(); // Check login status when the component is mounted
    }
};
</script>

<style scoped>
.auth-form {
    max-width: 400px;
    margin: 0 auto;
    padding: 20px;
    border-radius: 8px;
    background: #1e1e1e;
    /* Dark background for the form */
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.5);
    color: #e0e0e0;
    /* Light text color for contrast */
}

.auth-form h2 {
    font-size: 24px;
    font-weight: 600;
    margin-bottom: 20px;
    text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.3);
    /* Subtle shadow to enhance text readability */
}

.auth-form div {
    margin-bottom: 15px;
}

.auth-form label {
    display: block;
    font-size: 16px;
    font-weight: 600;
    margin-bottom: 5px;
}

.auth-form input {
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

.auth-form input::placeholder {
    color: #888;
}

.auth-form input:focus {
    border-color: #4caf50;
    box-shadow: 0 0 4px rgba(0, 123, 255, 0.2);
    outline: none;
}

.auth-form button {
    width: 100%;
    padding: 10px;
    background-color: #007bff;
    color: #fff;
    border: none;
    border-radius: 5px;
    cursor: pointer;
    font-size: 16px;
    transition: background-color 0.3s ease, transform 0.3s ease;
    margin-top: 10px;
}

.auth-form button:hover {
    background-color: #0056b3;
    transform: scale(1.05);
}

.auth-form button:last-of-type {
    background-color: #333;
    color: #fff;
    margin-top: 15px;
}

.auth-form button:last-of-type:hover {
    background-color: #f44336;
}
</style>