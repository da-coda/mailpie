import Vue from 'vue'
import vuetify from "@/plugins/vuetify";
import App from './App.vue'
import store from './store'
import router from './router'

Vue.config.productionTip = false
new Vue({
  vuetify,
  router,
  store,
  render: h => h(App)
}).$mount("#app");

