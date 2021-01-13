// import Vue from 'vue'
// import VueI18n from 'vue-i18n';
// import IdleVue from 'idle-vue';
// import StudentForm from './components/StudentForm'
// import ProgressWidget from "./components/ProgressWidget";
// import GoForth from "./components/GoForth";
// import MenuWidget from "./components/MenuWidget";
// import LangRu from './i18n/ru';
// import LangEn from './i18n/en';
//
// const eventsHub = new Vue();
// Vue.use(IdleVue, {
//   eventEmitter: eventsHub,
//   idleTime: 5*60*1000 // 5 minutes
// });
// Vue.use(VueI18n)
//
// const messages = {ru: LangRu.ru, en: LangEn.en};
//
// const i18n = new VueI18n({
//   locale: process.env.LANG,
//   fallbackLocale: 'en',
//   messages:  messages
// });
//
// Vue.component('studentForm', StudentForm)
// Vue.component('progressWidget', ProgressWidget)
// Vue.component('goForth', GoForth)
// Vue.component('menuWidget', MenuWidget)
//
// /* eslint-disable no-new */
// if(document.getElementById('app'))
//   new Vue({el: '#app', i18n})
// if(document.getElementById('progress-widget'))
//   new Vue({el: '#progress-widget', i18n})
// if(document.getElementById('go-forth-button'))
//   new Vue({el: '#go-forth-button', i18n})
// if(document.getElementById('menu-widget'))
//   new Vue({el: '#menu-widget', i18n})
