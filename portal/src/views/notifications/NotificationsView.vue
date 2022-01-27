<template>
    <StandardLayout>
        <div class="main-panel">
            <h1>{{ $t("notifications.title") }}</h1>

            <div class="notifications-container">
                <div class="notifications-filters">
                    <span class="active">{{ $t("notifications.filters.all") }}</span>
                    <span>{{ $t("notifications.filters.unread") }}</span>
                    <span>{{ $t("notifications.filters.comment") }}</span>
                    <span>{{ $t("notifications.filters.reply") }}</span>
                    <span>{{ $t("notifications.filters.mention") }}</span>
                </div>
                <NotificationsList></NotificationsList>
            </div>
        </div>
    </StandardLayout>
</template>

<script lang="ts">
import Vue from "vue";
import { mapState, mapGetters } from "vuex";
import StandardLayout from "../StandardLayout.vue";
import { GlobalState } from "@/store/modules/global";
import NotificationsList from "@/views/notifications/NotificationsList.vue";

export default Vue.extend({
    name: "NotificationsView",
    components: {
        NotificationsList,
        StandardLayout,
    },
    data() {
        return {};
    },
    computed: {
        ...mapState({ user: (s: GlobalState) => s.user.user }),
    },
    methods: {},
});
</script>

<style scoped lang="scss">
@import "../../scss/layout";
@import "../../scss/variables";

.main-panel h1 {
    font-size: 36px;
    margin: 0;
}

::v-deep .notifications {
    &-container {
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

    &-item {
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
</style>
