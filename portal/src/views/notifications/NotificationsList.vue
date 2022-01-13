<template>
    <ul>
        <li class="notifications-item notifications-item--reply">
            <div class="notifications-avatar">
                <UserPhoto :user="user" />
                <span class="notifications-item-icon"></span>
            </div>
            <div>
                Harriet Alexandra replied to your data view commment Harriet Alexandra replied to your data view commment
                <div class="notifications-timestamp">
                    1 day ago
                </div>
            </div>
            <i class="icon-ellipsis"></i>
        </li>
        <li class="notifications-item notifications-item--comment">
            <div class="notifications-avatar">
                <UserPhoto :user="user" />
                <span class="notifications-item-icon"></span>
            </div>
            <div>
                Harriet Alexandra replied to your data view commment
                <div class="notifications-timestamp">
                    1 day ago
                </div>
            </div>
            <i class="icon-ellipsis"></i>
        </li>
        <li class="notifications-item notifications-item--mention">
            <div class="notifications-avatar">
                <UserPhoto :user="user" />
                <span class="notifications-item-icon"></span>
            </div>
            <div>
                Harriet Alexandra replied to your data view commment
                <div class="notifications-timestamp">
                    1 day ago
                </div>
            </div>
            <i class="icon-ellipsis"></i>
        </li>
        <li class="notifications-item notifications-item--mention">
            <div class="notifications-avatar">
                <UserPhoto :user="user" />
                <span class="notifications-item-icon"></span>
            </div>
            <div>
                Harriet Alexandra replied to your data view commment
                <div class="notifications-timestamp">
                    1 day ago
                </div>
            </div>
            <i class="icon-ellipsis"></i>
        </li>
        <li class="notifications-item notifications-item--mention">
            <div class="notifications-avatar">
                <UserPhoto :user="user" />
                <span class="notifications-item-icon"></span>
            </div>
            <div>
                Harriet Alexandra replied to your data view commment
                <div class="notifications-timestamp">
                    1 day ago
                </div>
            </div>
            <i class="icon-ellipsis"></i>
        </li>
        <li class="notifications-item notifications-item--mention">
            <div class="notifications-avatar">
                <UserPhoto :user="user" />
                <span class="notifications-item-icon"></span>
            </div>
            <div>
                Harriet Alexandra replied to your data view commment
                <div class="notifications-timestamp">
                    1 day ago
                </div>
            </div>
            <i class="icon-ellipsis"></i>
        </li>
        <li class="notifications-item notifications-item--mention">
            <div class="notifications-avatar">
                <UserPhoto :user="user" />
                <span class="notifications-item-icon"></span>
            </div>
            <div>
                Harriet Alexandra replied to your data view commment
                <div class="notifications-timestamp">
                    1 day ago
                </div>
            </div>
            <i class="icon-ellipsis"></i>
        </li>
        <li class="notifications-item notifications-item--mention">
            <div class="notifications-avatar">
                <UserPhoto :user="user" />
                <span class="notifications-item-icon"></span>
            </div>
            <div>
                Harriet Alexandra replied to your data view commment
                <div class="notifications-timestamp">
                    1 day ago
                </div>
            </div>
            <i class="icon-ellipsis"></i>
        </li>
    </ul>
</template>

<script lang="ts">
import Vue from "vue";
import { mapState } from "vuex";
import { GlobalState } from "@/store/modules/global";
import UserPhoto from "../shared/UserPhoto.vue";

export default Vue.extend({
    name: "NotificationsList",
    components: {
        UserPhoto,
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

<style lang="scss" scoped>
@import "../../scss/mixins";

.notifications {
    &-item {
        @include flex(flex-start);
        color: #6a6d71;
        font-size: 14px;
        font-family: $font-family-light;
        margin-bottom: 10px;
        line-height: 1.4em;

        img {
            width: 35px;
            height: 35px;
            margin-right: 7px;
            margin-top: 0;
        }

        > div {
            padding-right: 7px;
        }

        &-icon {
            @include position(absolute, 24px 0 null 18px);
            display: block;
            border-radius: 50%;
            width: 17px;
            height: 17px;

            .notifications-item--reply & {
                background: url(/portal/src/assets/icon-reply-white.svg) no-repeat center center #ce596b;
                background-size: 11px;
            }

            .notifications-item--comment & {
                background: url(/portal/src/assets/icon-comment.svg) no-repeat center center #5268cc;
            }

            .notifications-item--mention & {
                background: url(/portal/src/assets/icon-mention.svg) no-repeat center center #52b5e4;
            }
        }
    }

    &-timestamp {
        font-size: 12px;
        margin-top: 3px;
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
</style>
