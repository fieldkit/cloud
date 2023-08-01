import Vue from "vue";
import { GlobalState, GlobalGetters } from "@/store/modules/global";
import { Services } from "@/api";

declare module "vue/types/vue" {
    interface Vue {
        $config: {
            baseUrl: string;
            sso: boolean;
        };
        $confirm: any;
        $loadAsset(path: string): any;
        $getters: GlobalGetters;
        $state: GlobalState;
        $services: Services;
        $seriousError: (error: Error) => void;
        $v: any;
    }
}

declare module "vue/types/options" {
    interface ComponentOptions<V extends Vue> {
      metaInfo?: any;
      validations?: any;
    }
  }