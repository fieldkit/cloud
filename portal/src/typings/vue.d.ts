import Vue from "vue";
import { Store } from "vuex";
import { GlobalState, GlobalGetters } from "../store/modules/global";

abstract class StrongVueClass extends Vue {
    public $store: Store<GlobalState>;


declare module "vue/types/vue" {
    interface Vue {
        $config: {
            baseUrl: string;
        };
        $loadAssets(path: string);
		$getters: GlobalGetters;
		$state: GlobalState;
    }
}
