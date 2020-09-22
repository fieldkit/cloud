<template>
    <div id="white-header" class="header">
        <div class="header-inner-section">
            <a class="menu-icon-container" v-on:click="toggleSidebar" v-bind:class="{ active: !isMenuNarrow }">
                <img alt="Menu icon" src="@/assets/icon-menu.svg" width="32" height="22"/>
            </a>
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
    props: {
        isMenuNarrow: {
            type: Boolean,
            default: false,
        },
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
        },
        toggleSidebar() {
            this.$emit("toggled");
        },
    },

});
</script>

<style scoped lang="scss">
@import '../../scss/mixins';
#white-header {
    width: auto;
    color: gray;
    text-align: right;
    overflow: hidden;
}
.header {

    &-inner-section {
        width: 100%;
        height: 69px;
        float: left;
        padding: 0 10px;
        box-sizing: border-box;
        @include flex(center, space-between);

        @include bp-down($sm) {
            padding: 0 20px;
        }

        @include bp-down($xs) {
            padding: 0 10px;
        }
    }

    &-account {
        padding-right: 62px;
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

.menu-icon-container {
    float: left;
    transition: all 0.33s;
    cursor: pointer;

    @include bp-down($md) {
        &.active {
            transform: translateX(175px);
        }
    }
}

.log-out,
.log-in {
    padding: 0 0 0 0;
    cursor: pointer;
}

</style>
