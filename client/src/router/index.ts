import Vue from 'vue'
import VueRouter from 'vue-router'

import ManageView from '@/views/ManageView.vue'
import NotFound from '@/views/404.vue'

import manageRoutes from './manageRoutes'
import Notebook from '@/views/Code/Notebook.vue'

Vue.use(VueRouter)

const routes = [
    {
        path: '/',
        name: 'Home',
        component: ManageView,
        redirect: manageRoutes[0].path,
        children: manageRoutes,
        meta: {
            alias: '首页'
        }
    }, {
        path: '*',
        name: '404',
        component: NotFound
    }
]

const router = new VueRouter({
    mode: 'history',
    base: process.env.BASE_URL,
    routes
})

export default router
