<template>
    <div class="terms">
        <img class="terms-header-logo" :alt="$t('layout.logo.alt')" src="@/assets/FieldKit_Logo_White.png" />
        <div class="terms-content">
            <img alt="Close" src="@/assets/icon-close.svg" class="close-button" v-on:click="goBack" />
            <div v-html="$t('terms.html')" class="terms-content-html"></div>
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
    background-image: linear-gradient(rgba(82, 181, 228, 0.6), rgba(27, 128, 201, 0.6));
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
        width: 80%;
        max-width: 800px;
        background-color: white;
        padding: 25px 45px 45px 45px;
        box-sizing: border-box;

        @include bp-down($sm) {
            width: calc(100% - 40px);
            padding: 22px 20px;
        }

        @include bp-down($xs) {
            width: calc(100% - 20px);
            padding: 22px 14px;
        }

        &-html {
            margin-top: 30px;
            p {
                text-align: justify;
            }
        }
    }

    &-body {
        font-size: 14px;
        line-height: 1.5;
    }
}

::v-deep h2 {
    font-size: 26px;
}
::v-deep h4 {
    font-size: 18px;
}
::v-deep h5 {
    font-size: 16px;
}
::v-deep {
    font-size: 14px;
}

.close-button {
    cursor: pointer;
    float: right;
}
</style>
