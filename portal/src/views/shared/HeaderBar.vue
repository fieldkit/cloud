<template>
    <div class="header">
        <div class="header-account">
            <UserPhoto v-if="user" :user="user" />
            <div>
                <router-link v-if="user" :to="{ name: 'editUser' }">
                    {{ user.name }}
                </router-link>
                <a class="log-out" v-if="isAuthenticated" v-on:click="logout">Log out</a>
                <router-link :to="{ name: 'login', query: { redirect: $route.fullPath } }" class="log-in" v-if="!isAuthenticated">
                    Log in
                </router-link>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import FKApi from "@/api/api";
import { mapState, mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";
import CommonComponents from "@/views/shared";
import { GlobalState } from "@/store/modules/global";

export default Vue.extend({
    name: "HeaderBar",
    components: {
        ...CommonComponents,
    },
    computed: {
        ...mapGetters({ isAuthenticated: "isAuthenticated" }),
        ...mapState({ user: (s: GlobalState) => s.user.user }),
        userImage() {
            if (this.$store.state.user.user.photo) {
                return this.$config.baseUrl + this.$store.state.user.user.photo.url;
            }
            return null;
        },
    },
    methods: {
        logout() {
            return this.$store.dispatch(ActionTypes.LOGOUT).then(() => {
                return this.$router.push({ name: "login" });
            });
        }
    },

});
</script>

<style scoped lang="scss">
@import '../../scss/mixins';

.header {
    background: #fff;
    box-shadow: 0 1px 4px 0 rgba(0, 0, 0, 0.12);
    z-index: $z-index-top;
    width: 100%;
    height: 66px;
    float: left;
    padding: 0 10px;
    box-sizing: border-box;
    @include flex(center, flex-end);

    @include bp-down($sm) {
        padding: 0 20px;
    }

    @include bp-down($xs) {
        padding: 0 10px;
    }

    &-account {
        padding-right: 62px;
        text-align: right;
        @include flex(center);

        @include bp-down($lg) {
            padding-right: 0;
        }

        a {
            display: block;
            font-size: 16px;
            color: initial;

            &:not(.log-out) {
                font-weight: 500;
            }
        }
    }
}

.log-out,
.log-in {
    padding: 0 0 0 0;
    cursor: pointer;
}

</style>
