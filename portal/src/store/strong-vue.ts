import Vue, { VueConstructor } from "vue";
import { Store } from "vuex";
import { GlobalState } from "./modules/global";

abstract class StrongVueClass extends Vue {
    public $store: Store<GlobalState>;
}
const StrongVue = Vue as VueConstructor<StrongVueClass>;

export default StrongVue;
