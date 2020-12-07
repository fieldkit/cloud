<template>
    <div class="transfer-station">
        <UserPicker @select="onUser" />
        <button :disabled="!user" v-on:click="transfer" class="button">Transfer</button>
    </div>
</template>

<script lang="ts">
import Vue, { PropType } from "vue";
import CommonComponents from "@/views/shared";
import UserPicker from "./UserPicker.vue";
import FKApi, { EssentialStation, SimpleUser } from "@/api/api";

export default Vue.extend({
    name: "TransferStation",
    components: { ...CommonComponents, UserPicker },
    props: {
        station: {
            type: Object as PropType<EssentialStation>,
            required: true,
        },
    },
    data(): {
        user: SimpleUser | null;
    } {
        return {
            user: null,
        };
    },
    methods: {
        onUser(user: SimpleUser): void {
            this.user = user;
        },
        async transfer(): Promise<void> {
            if (this.user) {
                await new FKApi().adminTransferStation(this.station.id, this.user.id);
                this.$emit("transferred", this.user);
            }
        },
    },
});
</script>

<style scoped>
.container {
    display: flex;
    flex-direction: column;
    padding: 20px;
    max-width: 900px;
}
</style>
