<template>
    <div class="items">
        <button
            class="item"
            :class="{ 'is-selected': index === selectedIndex }"
            v-for="(item, index) in items"
            :key="item.id"
            @click="selectItem(index)"
        >
            {{ item.name }}
        </button>
    </div>
</template>
<script lang="ts">
import Vue, { PropType } from "vue";

export default Vue.extend({
    props: {
        items: {
            type: Array as PropType<{ id: number; mention: string }[]>,
            required: true,
        },
        command: {
            type: Function,
            required: true,
        },
    },
    data(): {
        selectedIndex: number;
    } {
        return {
            selectedIndex: 0,
        };
    },
    watch: {
        items(): void {
            this.selectedIndex = 0;
        },
    },
    methods: {
        onKeyDown({ event }): boolean {
            if (event.key === "ArrowUp") {
                this.upHandler();
                return true;
            }

            if (event.key === "ArrowDown") {
                this.downHandler();
                return true;
            }

            if (event.key === "Enter") {
                console.log("mentions: enter-handler");
                this.enterHandler();
                return true;
            }

            if (event.key === "Tab") {
                console.log("mentions: tab-handler");
                this.enterHandler();
                return true;
            }

            return false;
        },
        upHandler(): void {
            this.selectedIndex = (this.selectedIndex + this.items.length - 1) % this.items.length;
        },
        downHandler(): void {
            this.selectedIndex = (this.selectedIndex + 1) % this.items.length;
        },
        enterHandler(): void {
            this.selectItem(this.selectedIndex);
        },
        selectItem(index: number): void {
            const item = this.items[index];
            console.log("mention", item);
            if (item) {
                this.command({ id: item.id, label: item.mention });
            }
        },
    },
});
</script>
<style lang="scss" scoped>
.items {
    position: relative;
    border-radius: 0.25rem;
    background: white;
    color: rgba(black, 0.8);
    overflow: hidden;
    font-size: 0.9rem;
    box-shadow: 0 0 0 1px rgba(0, 0, 0, 0.1), 0px 10px 20px rgba(0, 0, 0, 0.1);
}

.item {
    display: block;
    width: 100%;
    text-align: left;
    background: transparent;
    border: none;
    padding: 0.2rem 0.5rem;

    &.is-selected,
    &:hover {
        color: #a975ff;
        background: rgba(#a975ff, 0.1);
    }
}
</style>
