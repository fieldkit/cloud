<template>
    <div id="login-form-container">
        <img class="form-header-logo" alt="FieldKit Logo" src="../../assets/FieldKit_Logo_White.png" />
        <div v-if="!created">
            <form id="form" @submit.prevent="save">
                <h1>Create Your Account</h1>
                <div class="outer-input-container">
                    <div class="input-container">
                        <input
                            keyboardType="name"
                            autocorrect="false"
                            autocapitalizationType="none"
                            class="inputText"
                            v-model="form.name"
                        />
                        <span class="floating-label">Name</span>

                        <div class="validation-errors" v-if="$v.form.name.$error">
                            <div v-if="!$v.form.name.required">Name is a required field.</div>
                        </div>
                    </div>
                </div>
                <div class="outer-input-container">
                    <div class="input-container middle-container">
                        <input
                            keyboardType="email"
                            autocorrect="false"
                            autocapitalizationType="none"
                            class="inputText"
                            v-model="form.email"
                        />
                        <span class="floating-label">Email</span>

                        <div class="validation-errors" v-if="$v.form.email.$error">
                            <div v-if="!$v.form.email.required">Email is a required field.</div>
                            <div v-if="!$v.form.email.email">Must be a valid email address.</div>
                            <div v-if="!$v.form.email.taken">
                                This address appears to already be registered.
                                <router-link :to="{ name: 'recover' }" class="recover-link">
                                    Recover Account
                                </router-link>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="outer-input-container">
                    <div class="input-container middle-container">
                        <input name="password" secure="true" type="password" class="inputText" v-model="form.password" />
                        <span class="floating-label">Password</span>

                        <div class="validation-errors" v-if="$v.form.password.$error">
                            <div v-if="!$v.form.password.required">This is a required field.</div>
                            <div v-if="!$v.form.password.min">Password must be at least 10 characters.</div>
                        </div>
                    </div>
                </div>
                <div class="outer-input-container">
                    <div class="input-container middle-container">
                        <input
                            secure="true"
                            name="passwordConfirmation"
                            type="password"
                            class="inputText"
                            v-model="form.passwordConfirmation"
                        />
                        <span class="floating-label">Confirm Password</span>

                        <div class="validation-errors" v-if="$v.form.passwordConfirmation.$error">
                            <div v-if="!$v.form.passwordConfirmation.required">Confirmation is a required field.</div>
                            <div v-if="!$v.form.passwordConfirmation.sameAsPassword">Passwords must match.</div>
                        </div>
                    </div>
                    <div class="policy-terms-container">
                        By creating an account you agree to our
                        <span class="bold">Privacy Policy</span>
                        and
                        <span class="bold">Terms of Use.</span>
                    </div>
                </div>
                <button class="form-save-btn" type="submit">Create Account</button>
                <div>
                    <router-link :to="{ name: 'login' }" class="create-link">
                        Back to Log In
                    </router-link>
                </div>
            </form>
        </div>
        <div v-if="created">
            <div id="notifications">
                <div v-if="!resending">
                    <img alt="Success" src="../../assets/Icon_Success.png" width="57px" />
                    <p class="success">Account Created</p>
                    <div class="notification-text">We sent you an account validation email.</div>
                    <button class="form-save-btn" v-on:click="resend">Resend Email</button>
                    <router-link :to="{ name: 'login' }" class="create-link">
                        Back to Log In
                    </router-link>
                </div>
                <div v-if="resending">
                    <img alt="Resending" src="../../assets/Icon_Syncing2.png" width="57px" />
                    <p>Resending</p>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { required, email, minLength, sameAs } from "vuelidate/lib/validators";
import FKApi from "@/api/api";

export default {
    name: "CreateAccountView",
    data() {
        return {
            form: {
                name: "",
                email: "",
                password: "",
                passwordConfirmation: "",
            },
            available: true,
            creating: true,
            created: null,
            resending: false,
        };
    },
    validations: {
        form: {
            name: {
                required,
            },
            email: {
                required,
                email,
                taken: function () {
                    return this.available;
                },
            },
            password: { required, min: minLength(10) },
            passwordConfirmation: { required, min: minLength(10), sameAsPassword: sameAs("password") },
        },
    },
    methods: {
        save() {
            this.$v.form.$touch();
            if (this.$v.form.$pending || this.$v.form.$error) {
                return;
            }

            this.creating = true;

            return new FKApi()
                .register(this.form)
                .then((created) => {
                    this.created = created;
                })
                .catch((error) => {
                    if (error.status === 400) {
                        this.available = false;
                    } else {
                        return this.$seriousError(error);
                    }
                })
                .finally(() => {
                    this.creating = false;
                });
        },
        resend() {
            this.resending = true;
            return new FKApi().resendCreateAccount(this.created.id).then(() => {
                this.resending = false;
            });
        },
    },
};
</script>

<style scoped>
#login-form-container {
    width: 100%;
    min-height: 100%;
    background-image: linear-gradient(#52b5e4, #1b80c9);
}
#form,
#notifications {
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
.validation-errors {
    color: #c42c44;
    display: block;
    font-size: 14px;
}
.recover-link {
    color: black;
    font-weight: bold;
    margin-top: 10px;
    margin-bottom: 30px;
    display: block;
}
</style>
