<template>
    <div class="user-picker">
        <multiselect
            v-model="selected"
            label="name"
            track-by="id"
            placeholder="User's name or email"
            open-direction="bottom"
            :options="users"
            :multiple="false"
            :searchable="true"
            :loading="busy"
            :internal-search="false"
            :clear-on-select="false"
            :close-on-select="true"
            :show-no-results="true"
            :hide-selected="false"
            @search-change="search"
            @input="select"
        >
            <template slot="singleLabel" slot-scope="props">
                <div class="option__desc">
                    <span class="option__title">{{ props.option.name }} ({{ props.option.email }})</span>
                </div>
            </template>
            <template slot="option" slot-scope="props">
                <div class="option__desc">
                    <span class="option__title">{{ props.option.name }} ({{ props.option.email }})</span>
                </div>
            </template>
            <span slot="noResult">Oops! No users found.</span>
        </multiselect>
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import CommonComponents from "@/views/shared";
import FKApi, { SimpleUser } from "@/api/api";

import "vue-multiselect/dist/vue-multiselect.min.css";

export default Vue.extend({
    name: "UserPicker",
    components: { ...CommonComponents },
    props: {},
    data(): {
        busy: boolean;
        users: SimpleUser[];
        selected: SimpleUser[];
    } {
        return {
            busy: false,
            users: [],
            selected: [],
        };
    },
    methods: {
        async search(query: string): Promise<void> {
            this.busy = true;
            if (query.length == 0) {
                this.users = [];
            } else {
                const response = await new FKApi().adminSearchUsers(query);
                this.users = response.users;
            }
            this.busy = false;
        },
        select(user: SimpleUser): void {
            this.$emit("select", user);
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
