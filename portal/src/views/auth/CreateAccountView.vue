<template>
    <div id="form-container">
        <img class="form-header-logo" alt="FieldKit Logo" src="../../assets/FieldKit_Logo_White.png" />
        <div>
            <div id="form">
                <h1>Create Your Account</h1>
                <div class="outer-input-container">
                    <div class="input-container">
                        <input keyboardType="name" autocorrect="false" autocapitalizationType="none" required="" class="inputText" />
                        <span class="floating-label">Name</span>
                    </div>
                </div>
                <div class="outer-input-container">
                    <div class="input-container middle-container">
                        <input required="" keyboardType="email" autocorrect="false" autocapitalizationType="none" class="inputText" />
                        <span class="floating-label">Email</span>
                    </div>
                </div>
                <div class="outer-input-container">
                    <div class="input-container middle-container">
                        <input name="password" required="" secure="true" type="password" class="inputText" />
                        <span class="floating-label">Password</span>
                    </div>
                </div>
                <div class="outer-input-container">
                    <div class="input-container middle-container">
                        <input name="confirmPassword" required="" secure="true" ref="confirmPassword" type="password" class="inputText" />
                        <span class="floating-label">Confirm Password</span>
                    </div>
                    <div class="policy-terms-container">
                        By creating an account you agree to our
                        <span class="bold">Privacy Policy</span>
                        and
                        <span class="bold">Terms of Use.</span>
                    </div>
                </div>
                <button class="form-save-btn">Create Account</button>
                <div>
                    <router-link :to="{ name: 'login' }" class="create-link">
                        Back to Log In
                    </router-link>
                </div>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";

@Component({
    components: {},
})
export default class CreateAccountView extends Vue {}
</script>

<style scoped>
#form-container {
    width: 100%;
    min-height: 100%;
    background-image: linear-gradient(#52b5e4, #1b80c9);
}
#form {
    width: 460px;
    background-color: white;
    display: inline-block;
    text-align: center;
    padding-bottom: 60px;
    padding-top: 20px;
}
#notification {
    padding-top: 58px;
}
#reset-container {
    padding-top: 78px;
}
h1 {
    font-weight: 500;
    font-size: 24px;
    margin-bottom: 55px;
}
.bold {
    font-weight: bold;
}
.form-save-btn {
    margin-top: 50px;
    width: 300px;
    height: 45px;
    background-color: #ce596b;
    border: none;
    color: white;
    font-size: 18px;
    font-weight: 600;
    border-radius: 5px;
}
.disabled {
    margin-top: 20px;
    width: 300px;
    height: 50px;
}
.outer-input-container {
    width: 300px;
    margin: auto;
}
.input-container {
    margin: auto;
    width: 100%;
    text-align: left;
}
.middle-container {
    margin-top: 22px;
}
input {
    background: none;
    border: 0;
    border-bottom: 2px solid #d8dce0;
    outline: 0;
    font-size: 18px;
    padding-bottom: 2px;
}
ul {
    list-style-type: none;
    padding: 0;
}
li {
    display: inline-block;
    margin: 0 10px;
}
.policy-terms-container {
    font-size: 13px;
    line-height: 16px;
    text-align: left;
}
.create-link {
    cursor: pointer;
    display: block;
    margin-top: 30px;
    font-size: 14px;
    font-weight: bold;
}
.form-header-logo {
    width: 210px;
    margin-top: 150px;
    margin-bottom: 86px;
}
</style>
