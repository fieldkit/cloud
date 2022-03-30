<template>
    <div id="app">
        <router-view />
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import * as ActionTypes from "@/store/actions";
import { AuthenticationRequiredError } from "@/api";
import { getPartnerCustomization, PartnerCustomization } from "./views/shared/partners";

export default Vue.extend({
    async beforeMount(): Promise<void> {
        try {
            this.applyCustomClasses();
            await this.$store.dispatch(ActionTypes.INITIALIZE);
        } catch (err) {
            console.log("initialize error", err, err.stack);
        }
    },
    mounted(): void {
        this.setCustomFavicon();
        this.setCustomPageTitle();
    },
    computed: {
        partnerCustomization(): PartnerCustomization | null {
            return getPartnerCustomization();
        },
    },
    beforeUpdate(): void {
        this.applyCustomClasses();
    },
    errorCaptured(err): boolean {
        console.log("vuejs:error-captured", JSON.stringify(err));
        if (AuthenticationRequiredError.isInstance(err)) {
            this.$router.push({ name: "login", query: { after: this.$route.path } });
            return false;
        }
        return true;
    },
    methods: {
        applyCustomClasses(): void {
            if (this.partnerCustomization != null) {
                document.body.classList.add(this.partnerCustomization.class);
            }
        },
        setCustomFavicon(): void {
            const faviconEl = document.getElementById("favicon") as HTMLAnchorElement;
            if (this.partnerCustomization != null) {
                faviconEl.href = window.location.origin + this.partnerCustomization.icon;
            }
        },
        setCustomPageTitle(): void {
            if (this.partnerCustomization != null) {
                document.title = this.partnerCustomization.title;
            }
        },
    },
});
</script>
<style lang="scss">
@import "scss/mixins";
@import "scss/typography";
@import "scss/icons";
@import "icomoon/style.css";

html {
}

html,
body,
#app {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
}

body {
    --color-primary: #{$color-fieldkit-primary};
    --color-secondary: #{$color-fieldkit-secondary};
    --color-dark: #{$color-fieldkit-dark};
    --color-border: #{$color-fieldkit-border};
    --color-danger: #{$color-fieldkit-danger};
    --font-family-medium: #{$font-family-fieldkit-medium};
    --font-family-light: #{$font-family-fieldkit-light};
    --font-family-bold: #{$font-family-fieldkit-bold};

    text-align: center;
    margin: 0;
    padding: 0;
    flex-shrink: 0;
    color: var(--color-dark);
    font-family: var(--font-family-medium), Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;

    &.floodnet {
        --color-primary: #{$color-floodnet-primary};
        --color-secondary: #{$color-floodnet-dark};
        --color-dark: #{$color-floodnet-dark};
        --color-border: #{$color-floodnet-border};
        --color-danger: #{$color-fieldkit-danger};
        --font-family-medium: #{$font-family-floodnet-medium};
        --font-family-light: #{$font-family-floodnet-medium};
        --font-family-bold: #{$font-family-floodnet-bold};
        // v-calendar overwrite
        --blue-600: #{$color-floodnet-primary};
        --blue-200: #{lighten($color-floodnet-primary, 11%)};
    }
}

body:not(.disable-scrolling) {
    overflow-y: scroll;
}

body.disable-scrolling {
    margin-right: 14px; /* We need width of the scrollbars! */
}

body.blue-background {
    background-color: #1b80c9;

    @include bp-down($md) {
        background-color: #fff;
    }

    &.floodnet {
        @include bp-up($md) {
            background-color: var(--color-dark);
        }
    }
}

html.map-view {
    height: 100%;
}

body.map-view {
    height: 100%;
}

a {
    text-decoration: none;
    color: inherit;
}

button {
    cursor: pointer;
    color: inherit;

    body.floodnet & {
        font-family: $font-family-floodnet-button;
    }
}

.main-panel {
    width: auto;
    text-align: left;
    color: #2c3e50;
}

.main-panel h1 {
    font-size: 36px;
    margin-top: 40px;
}

h1 {
    font-family: var(--font-family-bold);
}

ul {
    margin: 0;
    padding: 0;
}

li {
    list-style-type: none;
}

.vue-treeselect__single-value {
    color: inherit;
}

.date-picker input {
    color: inherit;
}

.vue-treeselect__control {
    border: 1px solid var(--color-border);
}

.vc-nav-item {
    color: #fff !important;

    &.is-active {
        color: var(--color-dark) !important;
    }
}

.vc-nav-header * {
    color: #fff !important;
}
</style>
