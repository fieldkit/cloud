<template>
    <div class="header">
        <router-link :to="{ name: 'projects' }">
            <img :alt="$t('layout.logo.alt')" id="header-logo" src="@/assets/logo-fieldkit.svg" />
        </router-link>
        <div
            class="header-account"
            v-on:click="onAccountClick()"
            v-on:mouseenter="onAccountHover($event)"
            v-on:mouseleave="onAccountHover($event)"
        >
            <div class="header-avatar">
                <i class="badge">
                    <span>10</span>
                </i>
                <UserPhoto v-if="user" :user="user" />
                <span v-if="isAccountHovered" class="triangle"></span>
            </div>
            <a v-if="user" class="header-account-name">{{ firstName }}</a>
            <router-link :to="{ name: 'login', query: { redirect: $route.fullPath } }" class="log-in" v-if="!isAuthenticated">
                {{ $t("layout.header.login") }}
            </router-link>
            <div class="notifications-container" v-bind:class="{ active: isAccountHovered }">
                <header class="notifications-header">
                    <span>{{ $t("notifications.title") }}</span>
                    <div class="flex">
                        <router-link v-if="user && user.admin" :to="{ name: 'adminMain' }">
                            {{ $t("layout.header.admin") }}
                        </router-link>
                        <router-link v-if="user" :to="{ name: 'editUser' }" :title="$t('layout.header.myAccount')">
                            <img src="@/assets/icon-account.svg" alt="My Account" />
                        </router-link>
                        <a class="log-out" v-if="isAuthenticated" v-on:click="logout" :title="$t('layout.header.logout')">
                            <img src="@/assets/icon-logout.svg" alt="Logout" />
                        </a>
                    </div>
                </header>

                <NotificationsList></NotificationsList>

                <footer class="notifications-footer">
                    <button>{{ $t("notifications.viewAllButton") }}</button>
                    <button>{{ $t("notifications.dismissAllButton") }}</button>
                </footer>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import { mapState, mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";
import CommonComponents from "@/views/shared";
import NotificationsList from "@/views/notifications/NotificationsList.vue";
import { GlobalState } from "@/store/modules/global";

export default Vue.extend({
    name: "HeaderBar",
    components: {
        ...CommonComponents,
        NotificationsList,
    },
    data(): { isAccountHovered: boolean } {
        return {
            isAccountHovered: false,
        };
    },
    computed: {
        ...mapGetters({ isAuthenticated: "isAuthenticated" }),
        ...mapState({ user: (s: GlobalState) => s.user.user }),
        firstName(): string {
            if (!this.user) {
                return "";
            }
            return this.user.name.split(" ")[0];
        },
    },
    methods: {
        async logout(): Promise<void> {
            await this.$store.dispatch(ActionTypes.LOGOUT).then(() => {
                // Action handles where we go after.
            });
        },
        onAccountHover(event: Event): void {
            if (window.screen.availWidth < 768 && event.type == "mouseenter") {
                return;
            }
            this.isAccountHovered = !this.isAccountHovered;
        },
        onAccountClick(): void {
            this.isAccountHovered = !this.isAccountHovered;
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/mixins";
@import "../../scss/notifications";

.header {
    background: #fff;
    box-shadow: 0 1px 4px 0 rgba(0, 0, 0, 0.12);
    width: 100%;
    float: left;
    padding: 0 10px;
    box-sizing: border-box;
    z-index: $z-index-header;
    flex: 0 0 66px;
    @include flex(center, flex-end);

    @include bp-down($md) {
        padding: 0 10px;
        height: 54px;
        position: fixed;

        ::v-deep + * {
            margin-top: 54px;
        }
    }

    &-account {
        padding-right: 85px;
        text-align: right;
        position: relative;
        height: 100%;
        @include flex(center);

        @include bp-down($lg) {
            padding-right: 14px;
        }

        @include bp-down($sm) {
            padding-right: 0;
        }

        &-name {
            font-size: 16px;
            font-weight: 500;
            z-index: $z-index-top;

            @include bp-down($sm) {
                display: none;
            }
        }

        &:after {
            content: "";
            background: url("../../assets/icon-chevron-dropdown.svg") no-repeat center center;
            width: 10px;
            height: 10px;
            transition: all 0.33s;
            transform: translateY(-50%);
            cursor: pointer;
            @include position(absolute, 50% 69px null null);

            @include bp-down($lg) {
                right: 0;
            }

            @include bp-down($sm) {
                display: none;
            }
        }

        &:hover {
            &:after {
                transform: rotate(180deg) translateY(50%);
            }
        }
    }

    &-avatar {
        position: relative;
        cursor: pointer;
    }
}

.triangle {
    @include position(absolute, null null -10px 5px);
    z-index: $z-index-top;
    width: 0;
    height: 0;
    border-style: solid;
    border-width: 0 15px 12px 15px;
    border-color: transparent transparent #fff transparent;
    filter: drop-shadow(0px -2px 1px rgba(0, 0, 0, 0.1));

    @include bp-down($md) {
        left: 2px;
        border-width: 0 12px 9px 12px;
    }
}

.log-out,
.log-in {
    padding: 0 0 0 0;
    cursor: pointer;
}

::v-deep .default-user-icon {
    margin: 0 10px 0 0;

    @include bp-down($md) {
        width: 30px;
        height: 30px;
    }

    @include bp-down($sm) {
        margin: 0;
    }
}

#header-logo {
    @include position(fixed, 18px null null 50%);
    width: 103px;
    transform: translateX(-50%);

    @include bp-up($sm) {
        width: 140px;
    }

    @include bp-up($md) {
        display: none;
    }
}

.badge {
    @include position(absolute, -5px null null -7px);
    height: 20px;
    width: 20px;
    background: #1b80c9;
    border-radius: 50%;

    > * {
        @include position(absolute, 5px null null 50%);
        transform: translateX(-50%);
        color: #fff;
        font-size: 11px;
        font-style: normal;
        font-family: $font-family-bold;
    }
}

.flex {
    display: flex;
}

button {
    padding: 0;
    border: 0;
    outline: 0;
    box-shadow: none;
    cursor: pointer;
    background: transparent;
}

.notifications-container {
    background: #fff;
    transition: opacity 0.25s, max-height 0.33s;
    text-align: left;
    box-sizing: border-box;
    box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.5);
    border: solid 1px #e9e9e9;
    height: 80vh;
    flex-direction: column;
    width: 320px;
    z-index: -1;
    opacity: 0;
    visibility: hidden;
    @include flex();
    @include position(absolute, calc(100% + 1px) 70px null null);

    @include bp-down($lg) {
        top: 100%;
        right: 30px;
    }

    @include bp-down($sm) {
        right: -10px;
        height: calc(100vh - 55px);
    }

    @include bp-down($xs) {
        width: 100vw;
    }

    &.active {
        opacity: 1 !important;
        visibility: visible;
        z-index: initial;
    }

    a {
        padding: 8px 12px;
        font-size: 14px;
        display: block;
        user-select: none;
    }

    > ul {
        overflow-y: auto;
        padding: 10px;
    }
}
</style>
