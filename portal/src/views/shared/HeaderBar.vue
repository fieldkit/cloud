<template>
    <div class="header">
        <router-link :to="{ name: 'projects' }">
            <Logo />
        </router-link>
        <div
            class="header-account"
            :class="isAuthenticated ? 'loggedin' : ''"
            v-on:click="onAccountClick()"
            v-on:mouseenter="onAccountHover($event)"
            v-on:mouseleave="onAccountHover($event)"
        >
            <div v-if="user" class="header-avatar">
                <UserPhoto v-if="user" :user="user" />
                <span v-if="isAccountHovered" class="triangle"></span>
            </div>

            <a v-if="user" class="header-account-name">{{ firstName }}</a>

            <router-link :to="{ name: 'login', query: { redirect: $route.fullPath } }" class="log-in" v-if="!isAuthenticated">
                {{ $t("layout.header.login") }}
            </router-link>

            <div v-if="user" class="notifications-container" v-bind:class="{ active: isAccountHovered && !hiding }">
                <header class="notifications-header">
                    <span class="notifications-header-text">{{ $t("layout.header.myAccount") }}</span>
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

                <template v-if="false && numberOfUnseenNotifications > 0">
                    <NotificationsList v-on:notification-click="notificationNavigate"></NotificationsList>

                    <footer class="notifications-footer">
                        <button>{{ $t("notifications.viewAllButton") }}</button>
                        <button v-on:click="markAllSeen()">{{ $t("notifications.dismissAllButton") }}</button>
                    </footer>
                </template>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import { mapState, mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";
import { MarkNotificationsSeen } from "@/store";
import CommonComponents from "@/views/shared";
import NotificationsList from "@/views/notifications/NotificationsList.vue";
import { GlobalState } from "@/store/modules/global";
import Logo from "@/views/shared/Logo.vue";

export default Vue.extend({
    name: "HeaderBar",
    components: {
        ...CommonComponents,
        NotificationsList,
        Logo,
    },
    data(): { isAccountHovered: boolean; hiding: boolean } {
        return {
            isAccountHovered: false,
            hiding: false,
        };
    },
    computed: {
        ...mapGetters({
            isAuthenticated: "isAuthenticated",
            notifications: "notifications",
            numberOfUnseenNotifications: "numberOfUnseenNotifications",
        }),
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
            // console.log("hover", this.hiding, this.isAccountHovered);

            if (this.hiding && this.isAccountHovered) {
                this.isAccountHovered = false;
                this.hiding = false;
                return;
            }

            if (window.screen.availWidth < 768 && event.type == "mouseenter") {
                return;
            }
            this.isAccountHovered = !this.isAccountHovered;
        },
        onAccountClick(): void {
            this.hiding = true;
        },
        async markAllSeen(): Promise<void> {
            this.hiding = true;
            await this.$store.dispatch(new MarkNotificationsSeen([]));
        },
        notificationNavigate(ev: Event, notification: Notification): Promise<void> {
            console.log("notification", notification);
            return Promise.resolve();
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/mixins";

.header {
    background: #fff;
    box-shadow: 0 1px 4px 0 rgba(0, 0, 0, 0.12);
    width: 100%;
    float: left;
    padding: 0 10px;
    box-sizing: border-box;
    z-index: $z-index-header;
    flex: 0 0 65px;
    @include flex(center, flex-end);

    @include bp-down($md) {
        padding: 0 10px;
        height: 54px;
        position: fixed;

        ::v-deep + * {
            margin-top: 54px;
        }
    }

    > a {
        display: flex;
        align-items: center;
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
    }

    .loggedin {
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

    &-menu {
        overflow: hidden;
        background: #fff;
        transition: opacity 0.25s, max-height 0.33s;
        opacity: 0;
        visibility: hidden;
        text-align: left;
        min-width: 183px;
        box-sizing: border-box;
        box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.5);
        z-index: -1;
        @include position(absolute, calc(100% - 5px) 70px null null);

        @include bp-down($lg) {
            @include position(fixed, 60px 10px null unset);
        }

        &.active {
            opacity: 1 !important;
            visibility: visible;
            border: solid 1px #e9e9e9;
            z-index: initial;
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

.badge {
    @include position(absolute, -5px null null -7px);
    height: 20px;
    width: 20px;
    background: var(--color-primary);
    border-radius: 50%;

    > * {
        @include position(absolute, 5px null null 50%);
        transform: translateX(-50%);
        color: #fff;
        font-size: 11px;
        font-style: normal;
        font-family: $font-family-bold;

        body.floodnet & {
            color: var(--color-dark);
        }
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

.notifications {
    &-header {
        @include flex(center, space-between);
        height: 50px;
        /*
        border-bottom: solid 1px #d8dce0;
        margin-bottom: 15px;
        */
        letter-spacing: 0.1px;
        padding: 0 13px;

        &-text {
            font-size: 20px;
        }
    }

    &-footer {
        border-top: solid 1px #d8dce0;
        margin-top: auto;
        @include flex(center, space-between);

        button {
            padding: 13px 15px 10px 15px;
            font-size: 16px;
            font-weight: 900;
            color: #2c3e50;
        }
    }

    &-container {
        background: #fff;
        transition: opacity 0.25s, max-height 0.33s;
        text-align: left;
        box-sizing: border-box;
        box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.5);
        border: solid 1px #e9e9e9;
        /*height: 80vh;*/
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
}

#header-logo {
    display: none;

    @include bp-down($md) {
        @include position(fixed, null null null 50%);
        @include flex(center);
        font-size: 32px;
        height: 50px;
        transform: translateX(-50%);
    }

    @include bp-down($xs) {
        font-size: 26px;
    }
}
</style>
