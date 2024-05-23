import Vue from 'vue'
import App from './App.vue'


Vue.config.productionTip = false

import { Button, Empty, Popup } from 'vant';
Vue.component(Button.name, Button);
Vue.use(Empty);
Vue.use(Popup);

new Vue({
  render: h => h(App),
}).$mount('#app')
