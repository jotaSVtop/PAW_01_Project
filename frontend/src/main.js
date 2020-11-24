import Vue from 'vue'
import VueRouter from 'vue-router'
import VueGeolocation from 'vue-browser-geolocation'
import * as VueGoogleMaps from 'vue2-google-maps'
import axios from 'axios'
import VueAxios from 'vue-axios'
import BootstrapVue from 'bootstrap-vue'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import App from './App.vue'
import Login from './pages/Login.vue'
import Map from './pages/Map.vue'
import ListZones from './pages/ListZones.vue'
import Register from './pages/Register.vue'
import LocalControl from './pages/LocalControl.vue'
import AddRemoveZone from './pages/AddRemoveZone.vue'
import AdminPage from './pages/AdminPage.vue'
import VueJwtDecode from 'vue-jwt-decode'

Vue.config.productionTip = false
Vue.use(VueRouter)
Vue.use(VueAxios, axios)
Vue.use(BootstrapVue)
Vue.use(VueGeolocation)
Vue.use(VueJwtDecode)
Vue.use(VueGoogleMaps,  {
  load: {
    key: 'AIzaSyDvzBG1YC-EqqOau_4BMMAgK7p-t9nrYjE',
    autobindAllEvents: true,
  },
  installComponents: true,
})


const router = new VueRouter({
  routes: [
    {
      path: '/',
      name: 'login',
      component: Login
    },
    {
      path: '/map',
      name: 'map',
      component: Map
    },
    {
      path: '/zones',
      name: 'zones',
      component: ListZones
    },
    {
      path: '/admin/users',
      name: 'register',
      component: Register
    },
    {
      path: '/zones/:id',
      name: 'localControl',
      component:LocalControl
    },
    {
      path: '/admin/zones',
      name: 'adminZones',
      component: AddRemoveZone
    },
    {
      path: '/admin',
      name: 'admin',
      component: AdminPage 
    }
  ]
})

router.beforeEach((to, from, next) => {
  try{
    let jwtDecoded = VueJwtDecode.decode(localStorage.getItem('jwt'))
    if(to.name == 'register' && !jwtDecoded.admin){
      next({name: 'login'})
    }else if(to.name == 'adminZones' && !jwtDecoded.admin){
      next({name: 'login'})
    }else if(to.name == 'admin' && !jwtDecoded.admin){
      next({name: 'login'})
    }else{
      next()
    }
  }catch(e){
    if(to.name != 'login'){
      next({name: 'login'})
    }else{
      next()
    }
  }
})

new Vue({
  render: h => h(App),
  router
}).$mount('#app')
