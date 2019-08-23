<template>
    <div id="LoginForm">
        <h1>Log In to Your Account</h1>
        <p v-if="error">Your email or password is invalid</p>
        <input name="email" v-model="email" placeholder="Email" />
        <br />
        <br />
        <input name="password" type="password" placeholder="Password" v-model="password" />
        <br />
        <button v-on:click="login">Log In</button>
    </div>
</template>

<script>
import FKApi from "../../api/api";

export default {
    name: "LoginForm",
    props: {
        msg: String
    },
    data: () => {
        return { email: "", password: "", user: {}, error: false };
    },
    methods: {
        async login(event) {
            event.preventDefault();
            try {
                if (this.email == "" || this.password == "") {
                    this.error = true;
                }
                const api = new FKApi();
                const auth = await api.login(this.email, this.password);
                const isAuthenticated = await api.authenticated();
                if (isAuthenticated) {
                    this.userToken = auth;
                    this.error = false;
                    this.$router.push({ name: "dashboard" });
                } else {
                    this.error = true;
                }
            } catch (error) {
                this.error = false;
            }
        }
    }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1 {
    font-weight: lighter;
}
p {
    color: red;
}
button {
    margin-top: 50px;
    width: 70%;
    height: 50px;
    background-color: #ce596b;
    border: none;
    color: white;
    font-size: 18px;
    border-radius: 5px;
}
input {
    border: 0;
    border-bottom: 1px solid gray;
    outline: 0;
    height: 30px;
    width: 70%;
    font-size: 18px;
}
ul {
    list-style-type: none;
    padding: 0;
}
li {
    display: inline-block;
    margin: 0 10px;
}
div {
    width: 30%;
    background-color: white;
    display: inline-block;
    text-align: center;
    padding-bottom: 60px;
    padding-top: 20px;
}
</style>
