import Vue from "vue";
import { Store } from "vuex";
import { GlobalState, GlobalGetters } from "../store/modules/global";

declare module "vue/types/vue" {
    interface Vue {
        $config: {
            baseUrl: string;
        };
        $loadAssets(path: string);
        $getters: GlobalGetters;
        $state: GlobalState;
        $seriousError: (error: Error) => void;
    }
}
