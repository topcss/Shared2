import Vue from 'vue'
import App from './App.vue'


Vue.config.productionTip = false

import { Checkbox, Button, Empty, Popup } from 'vant';
Vue.component(Button.name, Button);
Vue.use(Empty);
Vue.use(Checkbox);
Vue.use(Popup);

new Vue({
  render: h => h(App),
}).$mount('#app')

const debounce = (fn, delay) => {
  let timer = null;
  return function () {
    let context = this;
    let args = arguments;
    clearTimeout(timer);
    timer = setTimeout(function () {
      fn.apply(context, args);
    }, delay);
  }
}

// 解决bug
const _ResizeObserver = window.ResizeObserver;
window.ResizeObserver = class ResizeObserver extends _ResizeObserver {
  constructor(callback) {
    callback = debounce(callback, 16);
    super(callback);
  }
}
