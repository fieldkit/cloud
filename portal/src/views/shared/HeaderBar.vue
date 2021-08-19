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
            <UserPhoto v-if="user" :user="user" />
            <a v-if="user" class="header-account-name">
                <b-badge>
                    {{ numberOfUnseenNotifications }}
                    <span class="sr-only" style="display: none">unread notifications.</span>
                </b-badge>
                {{ firstName }}
            </a>
            <router-link :to="{ name: 'login', query: { redirect: $route.fullPath } }" class="log-in" v-if="!isAuthenticated">
                {{ $t("layout.header.login") }}
            </router-link>
            <div class="header-account-menu" v-bind:class="{ active: isAccountHovered && !hiding }">
                <router-link v-if="user" :to="{ name: 'editUser' }">{{ $t("layout.header.myAccount") }}</router-link>
                <router-link v-if="user && user.admin" :to="{ name: 'adminMain' }">{{ $t("layout.header.admin") }}</router-link>
                <a class="log-out" v-if="isAuthenticated" v-on:click="logout">{{ $t("layout.header.logout") }}</a>

                <div v-if="numberOfUnseenNotifications > 0">
                    <div class="menu-item menu-heading">{{ $t("layout.header.notifications") }}</div>
                    <a
                        v-for="notification in notifications"
                        :key="notification.notificationId"
                        v-on:click="(ev) => notificationNavigate(ev, notification)"
                    >
                        {{ notificationDisplay(notification) }}
                    </a>
                    <a class="mark-all-as-seen" v-on:click="markAllSeen">
                        {{ $t("layout.header.markAllSeen") }}
                    </a>
                </div>
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
import { GlobalState } from "@/store/modules/global";

export default Vue.extend({
    name: "HeaderBar",
    components: {
        ...CommonComponents,
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
            console.log("hover", this.hiding, this.isAccountHovered);

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
        notificationDisplay(notification: Notification): string {
            return "Notification";
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
    height: 66px;
    float: left;
    padding: 0 10px;
    box-sizing: border-box;
    z-index: $z-index-header;
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

        * {
            font-weight: 500;
            color: #2c3e50;
            cursor: pointer;
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

            a,
            .menu-item {
                padding: 8px 17px;
                font-size: 14px;
                display: block;
                user-select: none;
                padding: 12px;
            }

            .menu-heading {
                font-weight: bold;
                background-color: #efefef;
            }
        }
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
        width: 25px;
        height: 25px;
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

.badge-secondary {
    color: #fff;
    background-color: #dc3545;
}

.badge {
    display: inline-block;
    padding: 0.5em 0.5em 0.5em 0.5em;
    font-size: 75%;
    font-weight: 700;
    line-height: 1;
    text-align: center;
    white-space: nowrap;
    vertical-align: baseline;
    border-radius: 0.25rem;
    transition: color 0.15s ease-in-out, background-color 0.15s ease-in-out, border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out;
}
</style>
