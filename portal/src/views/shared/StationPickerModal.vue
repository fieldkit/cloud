<template>
    <div class="modal-mask">
        <div class="modal-wrapper">
            <div class="modal-container">
                <StationPicker
                    :title="title"
                    :actionText="actionText"
                    :stations="stations"
                    :filter="filter"
                    :actionType="actionType"
                    @ctaClick="emitAction"
                    @close="$emit('close')"
                />
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import Vue, { PropType } from "vue";
import StationPicker from "./StationPicker.vue";
import { DisplayStation } from "@/store";

export default Vue.extend({
    name: "StationPickerModal",
    components: {
        StationPicker,
    },
    props: {
        title: {
            type: String,
            required: true,
        },
        actionType: {
            type: String,
            required: true,
        },
        actionText: {
            type: String,
            required: true,
        },
        stations: {
            type: Array as PropType<DisplayStation[]>,
            required: true,
        },
        filter: {
            type: Function as PropType<(station: DisplayStation) => boolean>,
            default: (station) => true,
        },
    },
    methods: {
        emitAction(ids): void {
            this.$emit(this.actionType, ids);
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/mixins";

.modal-mask {
    position: fixed;
    z-index: 9998;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    display: table;
    transition: opacity 0.3s ease;
}

.modal-wrapper {
    display: table-cell;
    vertical-align: middle;
}

.modal-container {
    box-sizing: border-box;
    width: 840px;
    max-width: 80%;
    margin: 0px auto;
    padding: 70px 45px 35px;
    background-color: #fff;
    box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.24);
    border: solid 1px #f4f5f7;
    border-radius: 2px;
    transition: all 0.3s ease;
    font-family: Helvetica, Arial, sans-serif;

    @include bp-down($md) {
        width: unset;
    }

    @include bp-down($sm) {
        width: unset;
        padding: 60px 30px 35px;
    }

    @include bp-down($xs) {
        max-width: 100%;
        padding: 40px 15px 35px;
    }
}

.modal-header h3 {
    margin-top: 0;
    color: #42b983;
}

.modal-body {
    margin: 20px 0;
}
/*
 * The following styles are auto-applied to elements with
 * transition="modal" when their visibility is toggled
 * by Vue.js.
 *
 * You can easily play with the modal transition by editing
 * these styles.
 */

.modal-enter {
    opacity: 0;
}

.modal-leave-active {
    opacity: 0;
}

.modal-enter .modal-container,
.modal-leave-active .modal-container {
    -webkit-transform: scale(1.1);
    transform: scale(1.1);
}
</style>
