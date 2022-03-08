<template>
    <ul>
        <li
            class="notifications-item"
            :class="`notifications-item--${notification.kind}`"
            v-for="notification in notificationsVisible"
            :key="notification.notificationId"
            v-on:click="onClick(notification)"
        >
            <div class="notifications-avatar">
                <UserPhoto :user="notification.user" />
                <span class="notifications-item-icon"><span :class="`icon-icon-${notification.kind}`"></span></span>
            </div>
            <div>
                <span>{{ notification.user.name }}</span>
                <span class="notification-kind">{{ display(notification) }}</span>
                <div class="notifications-timestamp">
                    {{ formatTimestamp(notification.createdAt) }}
                </div>
            </div>
        </li>
    </ul>
</template>

<script lang="ts">
import Vue, { PropType } from "vue";
import { mapGetters, mapState } from "vuex";
import { GlobalState } from "@/store/modules/global";
import UserPhoto from "../shared/UserPhoto.vue";
import moment from "moment";
import { Notification } from "@/store/modules/notifications";

export default Vue.extend({
    name: "NotificationsList",
    components: {
        UserPhoto,
    },
    data(): {
        notificationsData: Notification[];
    } {
        return {
            notificationsData: [],
        };
    },
    computed: {
        ...mapGetters({
            notifications: "notifications",
        }),
        ...mapState({ user: (s: GlobalState) => s.user.user }),
        notificationsVisible() {
            return this.notificationsList ?? this.notifications;
        },
    },
    props: {
        notificationsList: {
            type: Array as PropType<Notification[]>,
            required: false,
        },
    },
    methods: {
        display(notification: Notification): string {
            switch (notification.kind) {
                case "reply":
                    return "replied to your comment";
                case "comment":
                    return "commented on your data";
                case "mention":
                    return "mentioned you in a comment";
            }

            return "Notification";
        },
        onClick(notification: Notification): void {
            this.$emit("notification-click", notification);
        },
        formatTimestamp(timestamp: number): string {
            return moment(timestamp).fromNow();
        },
    },
});
</script>

<style lang="scss" scoped>
@import "../../scss/mixins";

.notifications {
    &-item {
        @include flex(center);
        color: #6a6d71;
        font-size: 14px;
        font-family: $font-family-light;
        margin-bottom: 10px;
        line-height: 1.4em;

        > div {
            padding-right: 7px;
        }

        &-icon {
            @include position(absolute, 24px 0 null 18px);
            display: block;
            border-radius: 50%;
            width: 17px;
            height: 17px;
            font-size: 10px;
            line-height: 18px;

            @include bp-down($md) {
                @include position(absolute, 26px 0 null 20px);
                width: 15px;
                height: 15px;
                font-size: 8px;
                line-height: 16px;
            }

            .notifications-item--reply & {
                background-color: #ce596b;
            }

            .notifications-item--comment & {
                background-color: #5268cc;
            }

            .notifications-item--mention & {
                background-color: #52b5e4;
            }
        }

        .notification-kind {
            color: #6a6d71;

            &::before {
                content: " ";
            }
        }

        ::v-deep .default-user-icon {
            width: 35px;
            height: 35px;
            margin-right: 7px;
            margin-top: 0;
        }
    }

    &-timestamp {
        font-size: 12px;
        color: #6a6d71;
    }

    &-avatar {
        position: relative;
    }
}

.icon-ellipsis {
    display: block;
    cursor: pointer;
    margin-left: auto;

    &:after {
        @include flex(flex-end);
        content: "...";
        color: #2c3e50;
        height: 17px;
        font-size: 32px;
        font-family: $font-family-bold;
        letter-spacing: -1.5px;
    }
}

.notifications-avatar {
    text-align: center;

    .icon-icon-reply,
    .icon-icon-comment,
    .icon-icon-mention {
        font-size: 10px;
        line-height: 18px;
        @include bp-down($md) {
            font-size: 8px;
            line-height: 16px;
        }
    }
    .icon-icon-reply:before,
    .icon-icon-comment:before,
    .icon-icon-mention:before {
        color: white;
    }
}
</style>
