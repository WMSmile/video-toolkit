import { createRouter, createWebHistory } from 'vue-router'

import Home from '../views/Home.vue'
import About from '../views/About.vue'
import Settings from '../views/Settings.vue'
import Convert from '../views/Convert.vue'
import Tagger from '../views/Tagger.vue'
import Batch from '../views/Batch.vue'
import ASR from '../views/ASR.vue'



const routes = [
  { path: '/', name: 'Home', component: Home },
  { path: '/about', name: 'About', component: About },
  { path: '/settings', name: 'Settings', component: Settings },
  { path: '/convert', name: 'Convert', component: Convert },
  { path: '/tagger', name: 'Tagger', component: Tagger },
  { path: '/batch', name: 'Batch', component: Batch },
  { path: '/asr', name: 'ASR', component: ASR }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router