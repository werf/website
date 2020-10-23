import Vue from 'vue'
import IdleVue from 'idle-vue';
import StudentForm from './components/StudentForm'
import ProgressWidget from "./components/ProgressWidget";
import GoForth from "./components/GoForth";
import MenuWidget from "./components/MenuWidget";

const eventsHub = new Vue();
Vue.use(IdleVue, {
  eventEmitter: eventsHub,
  idleTime: 5*60*1000 // 5 minutes
});

Vue.component('studentForm', StudentForm)
Vue.component('progressWidget', ProgressWidget)
Vue.component('goForth', GoForth)
Vue.component('menuWidget', MenuWidget)

/* eslint-disable no-new */
if(document.getElementById('app'))
  new Vue({el: '#app'})
if(document.getElementById('progress-widget'))
  new Vue({el: '#progress-widget'})
if(document.getElementById('go-forth-button'))
  new Vue({el: '#go-forth-button'})
if(document.getElementById('menu-widget'))
  new Vue({el: '#menu-widget'})
