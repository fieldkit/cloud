<template class="container">
    <section class="container">
        <header> Notes & Comments </header>
    </section>
</template>

<script>
    export default {
        name: "Comments",
        components: {},
        props: {
            type: {
                required: true
            },
        }
    };
</script>

<style lang="scss">
    @import "../../scss/global";

    .hide {
        display: none;
    }

    .container {
        margin-top: 20px;
        background: #fff;
        border-radius: 1px;
        border: 1px solid $color-border;
        box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.05);
    }

    header {
        @include flex(center, space-between);
        height: 52px;
        padding: 0 20px;
        border-bottom: 1px solid $color-border;
        font-size: 20px;
        font-weight: 500;
    }


</style>
