<template>
    <div></div>
</template>
<script>
import Vue from "@/store/strong-vue";

export default Vue.extend({
    name: "Zoho",
    mounted() {
        if (process.env.NODE_ENV !== "production") {
            console.log("zoho:disabled");
            return;
        }

        window.ZohoHCAsap =
            window.ZohoHCAsap ||
            function (a, b) {
                global.ZohoHCAsap[a] = b;
            };

        (function () {
            const d = document;
            const s = d.createElement("script");
            s.type = "text/javascript";
            s.defer = true;
            s.src = "https://desk.zoho.com/portal/api/web/inapp/514122000000349001?orgId=718189455";
            d.getElementsByTagName("head")[0].appendChild(s);
        })();
    },
    methods: {},
});
</script>
