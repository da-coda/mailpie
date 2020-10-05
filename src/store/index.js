import Vue from "vue";
import Vuex from "vuex";

import mails from "./mails";

Vue.use(Vuex);

export default new Vuex.Store({
    modules: {
        mails
    }
});