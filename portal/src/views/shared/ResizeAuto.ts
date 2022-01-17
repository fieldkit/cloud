import Vue from "vue";

export const ResizeAuto = Vue.extend({
    name: "ResizeAuto",
    methods: {
        resize(ev) {
            ev.target.style.height = "auto";
            ev.target.style.height = `${ev.target.scrollHeight}px`;
            console.log("resize:", ev.target.scrollHeight, ev.target.style.height);
        },
    },
    mounted() {
        this.$nextTick(() => {
            // workaround since $el.scrollHeight does not have the correct value without timeout
            setTimeout(() => {
                const el: any = this.$el;
                el.setAttribute("style", "height:" + this.$el.scrollHeight + "px");
                console.log("resize-timeout:", this.$el.scrollHeight);
            }, 200);
        });
        this.$el.addEventListener("input", this.resize);
    },
    beforeDestroy() {
        this.$el.removeEventListener("input", this.resize);
    },
    render(this: any) {
        return this.$scopedSlots.default({
            resize: this.resize,
        });
    },
});
