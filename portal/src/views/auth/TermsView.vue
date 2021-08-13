<template>
    <div class="terms">
        <img class="terms-header-logo" :alt="$t('layout.logo.alt')" src="@/assets/FieldKit_Logo_White.png" />
        <div class="terms-content">
            <img alt="Close" src="@/assets/icon-close.svg" class="close-button" v-on:click="goBack" />
            <p class="terms-title">{{ $t("terms.title") }}</p>
            <p class="terms-body">{{ $t("terms.body") }}</p>
            <form v-if="tncOutdated" @submit.prevent="save">
                <button class="form-submit" type="submit">{{ $t("terms.agreeButton") }}</button>
            </form>
        </div>
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import CommonComponents from "@/views/shared";

import { mapGetters, mapState } from "vuex";
import { GlobalState } from "@/store";
import * as ActionTypes from "@/store/actions";

export default Vue.extend({
    name: "CreateAccountView",
    components: {
        ...CommonComponents,
    },
    computed: {
        ...mapGetters({ isAuthenticated: "isAuthenticated", isTncValid: "isTncValid" }),
        ...mapState({
            user: (s: GlobalState) => s.user.user,
        }),
        tncOutdated() {
            return this.isAuthenticated && !this.isTncValid;
        },
    },
    methods: {
        async save(): Promise<void> {
            if (this.user) {
                const payload = {
                    id: this.user?.id,
                    name: this.user?.name,
                    email: this.user?.email,
                    bio: this.user?.bio,
                    tncDate: this.user?.tncDate,
                    tncAccept: true,
                };
                await this.$services.api.accept(this.user.id);
                await this.$store.dispatch(ActionTypes.REFRESH_CURRENT_USER, {});
                await this.$router.push({ name: "projects" }).catch((e) => {
                    console.log(e);
                });
            }
        },
        goBack(): void {
            if (window.history.length) {
                this.$router.go(-1);
            } else {
                this.$router.push("/");
            }
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/forms";
@import "../../scss/global";
@import "../../scss/mixins";

.logo {
    width: 150px;
    margin-bottom: 50px;
}

.terms {
    background-image: linear-gradient(#52b5e4, #1b80c9);
    padding: 40px 0;
    flex-direction: column;
    box-sizing: border-box;
    min-height: 100vh;
    color: #2c3e50;
    text-align: left;
    @include flex(center, center);

    &-header-logo {
        width: 211px;
        margin-bottom: 45px;

        @include bp-down($xs) {
            width: 117px;
            margin-bottom: 20px;
        }
    }
    &-content {
        width: 60%;
        text-align: justify;
        background-color: white;
        padding: 25px 45px 45px 45px;
    }
    &-title {
        font-size: 36px;
        font-weight: 900;
    }

    &-body {
        font-size: 16px;
        line-height: 1.5;
    }
}

.close-button {
    cursor: pointer;
    float: right;
}
</style>
