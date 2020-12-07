// import Vue from "vue";
// import { Store } from "vuex";
import { GlobalState, GlobalGetters } from "@/store/modules/global";
import { Services } from "@/api";

declare module "vue/types/vue" {
    interface Vue {
        $config: {
            baseUrl: string;
        };
        $confirm: any;
        $loadAsset(path: string): any;
        $getters: GlobalGetters;
        $state: GlobalState;
        $services: Services;
        $seriousError: (error: Error) => void;
    }
}
