<template>
    <div id="login-container">
        <img
            v-bind:style="{
                width: '210px',
                marginTop: '150px',
                marginBottom: '86px',
            }"
            alt="Fieldkit Logo"
            src="../assets/FieldKit_Logo_White.png"
        />
        <br />
        <LoginForm />
    </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import LoginForm from "../components/LoginForm.vue";

@Component({
    components: {
        LoginForm,
    },
})
export default class Login extends Vue {}
</script>

<style scoped>
#login-container {
    width: 100%;
    min-height: 100%;
    background-image: linear-gradient(#52b5e4, #1b80c9);
}
</style>
