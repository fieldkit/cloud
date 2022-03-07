<template>
    <StandardLayout>
        <div class="main-panel">
            <h1>{{ $t("notifications.title") }}</h1>

            <div class="notifications-body">
                <div class="notifications-filters">
                    <span :class="{ active: this.filter === '' }" v-on:click="viewAll()">{{ $t("notifications.filters.all") }}</span>
                    <span :class="{ active: this.filter === 'comment' }" v-on:click="viewKind('comment')">
                        {{ $t("notifications.filters.comment") }}
                    </span>
                    <span :class="{ active: this.filter === 'reply' }" v-on:click="viewKind('reply')">
                        {{ $t("notifications.filters.reply") }}
                    </span>
                    <span :class="{ active: this.filter === 'mention' }" v-on:click="viewKind('mention')">
                        {{ $t("notifications.filters.mention") }}
                    </span>
                </div>
                <NotificationsList
                    v-on:notification-click="notificationNavigate"
                    :notificationsList="notificationsList"
                ></NotificationsList>
                <div class="no-notifications" v-if="notificationsList.length === 0">
                    <span>{{ $t("layout.header.noNotifications") }}</span>
                </div>
            </div>
        </div>
    </StandardLayout>
</template>

<script lang="ts">
import Vue, { PropType } from "vue";
import { mapState, mapGetters } from "vuex";
import StandardLayout from "../StandardLayout.vue";
import { GlobalState } from "@/store/modules/global";
import NotificationsList from "@/views/notifications/NotificationsList.vue";
import { Notification } from "@/store/modules/notifications";

export default Vue.extend({
    name: "NotificationsView",
    components: {
        NotificationsList,
        StandardLayout,
    },
    data() {
        return {
            filter: "",
        };
    },
    computed: {
        ...mapGetters({
            notifications: "notifications",
        }),
        ...mapState({ user: (s: GlobalState) => s.user.user }),
        notificationsList(): Notification[] {
            return this.filter ? this.notifications.filter((item) => item.kind === this.filter) : this.notifications;
        },
    },
    methods: {
        viewAll() {
            this.filter = "";
        },
        viewKind(kind: string) {
            this.filter = kind;
        },
        notificationNavigate(notification: Notification) {
            if (notification.projectId) {
                return this.$router
                    .push({
                        name: "viewProject",
                        params: { id: notification.projectId },
                        hash: `#comment-id-${notification.postId}`,
                    })
                    .catch((err) => {
                        return;
                    });
            }
            if (notification.bookmark) {
                return this.$router
                    .push({
                        name: "exploreBookmark",
                        params: { bookmark: notification.bookmark },
                        hash: `#comment-id-${notification.postId}`,
                    })
                    .catch((err) => {
                        return;
                    });
            }
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/layout";
@import "../../scss/variables";

.main-panel h1 {
    font-size: 36px;
    margin: 0;
}

.notifications {
    &-body {
        margin-top: 50px;
        background: #fff;
        padding: 24px 24px 112.5px 23px;
        border-radius: 1px;
        border: solid 1px $color-border;

        @include bp-down($xs) {
            margin-top: 20px;
            padding: 24px 10px 24px 10px;
        }

        > ul {
            border-top: 1px solid #ddd;
            margin-top: 23px;
        }
    }

    ::v-deep &-item {
        margin-bottom: 0;
        padding: 14px 0 12px;
        border-bottom: 1px solid #ddd;
    }

    &-timestamp {
        margin-top: 0;
    }

    &-filters {
        span {
            font-size: 12px;
            margin: 0 9px;
            position: relative;
            cursor: pointer;

            &.active {
                font-weight: 900;

                &:after {
                    content: "";
                    height: 2px;
                    background: #2c3e50;
                    width: 100%;
                    @include position(absolute, null null -2px 0);
                }
            }
        }
    }
}

.no-notifications {
    margin-top: 15px;
}
</style>
