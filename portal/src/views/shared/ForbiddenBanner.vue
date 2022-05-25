<template>
    <div class="message-container">
        <img src="@/assets/icon-warning-error.svg" alt="Error" width="15px" />
        <div>
            <div class="message-title" v-if="title">{{ title }}</div>
            <div class="message-subtitle" v-if="subtitle">{{ subtitle }}</div>
        </div>
    </div>
</template>

<script>
export default {
    name: "ForbiddenBanner",
    props: {
        title: {
            type: String,
            required: false,
        },
        subtitle: {
            type: String,
            required: false,
        },
    },
};
</script>

<style scoped lang="scss">
@import "src/scss/mixins";

.message {
    &-container {
        display: flex;
        width: 460px;
        background-color: #f4f5f7;
        box-sizing: border-box;
        padding: 15px 15px;
        max-width: calc(100vw - 20px);
        align-items: center;
        justify-content: center;
        margin-bottom: 80px;
        line-height: 1.5;

        > img {
            margin-right: 10px;
        }

        @include bp-down($xs) {
            width: 330px;
        }
    }

    &-title {
        font-size: 16px;
        text-align: left;
    }

    &-subtitle {
        font-size: 12px;
        text-align: left;
    }
}
</style>
