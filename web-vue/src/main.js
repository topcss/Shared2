import Vue from 'vue'
import App from './App.vue'


Vue.config.productionTip = false

import { Button, Empty } from 'vant';
Vue.component(Button.name, Button);

Vue.use(Empty);

new Vue({
  render: h => h(App),
}).$mount('#app')
