import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

/* Layout */
import Layout from '@/layout'

/**
 * Note: sub-menu only appear when route children.length >= 1
 * Detail see: https://panjiachen.github.io/vue-element-admin-site/guide/essentials/router-and-nav.html
 *
 * hidden: true                   if set true, item will not show in the sidebar(default is false)
 * alwaysShow: true               if set true, will always show the root menu
 *                                if not set alwaysShow, when item has more than one children route,
 *                                it will becomes nested mode, otherwise not show the root menu
 * redirect: noRedirect           if set noRedirect will no redirect in the breadcrumb
 * name:'router-name'             the name is used by <keep-alive> (must set!!!)
 * meta : {
    roles: ['admin','editor']    control the page roles (you can set multiple roles)
    title: 'title'               the name show in sidebar and breadcrumb (recommend set)
    icon: 'svg-name'             the icon show in the sidebar
    breadcrumb: false            if set false, the item will hidden in breadcrumb(default is true)
    activeMenu: '/example/list'  if set path, the sidebar will highlight the path you set
  }
 */

/**
 * constantRoutes
 * a base page that does not have permission requirements
 * all roles can be accessed
 */
export const constantRoutes = [
  {
    path: '/login',
    component: () => import('@/views/login/index'),
    hidden: true
  },

  {
    path: '/404',
    component: () => import('@/views/404'),
    hidden: true
  },

  {
    path: '/',
    component: Layout,  
    redirect: '/task',
    // hidden: true,
    children: [{
      path: 'dashboard',
      name: 'Dashboard',
      component: () => import('@/views/dashboard/index'),
      meta: { title: '首页', icon: 'dashboard' }
    }]
  },
  {
    path: '/task',
    component: Layout,
    children: [
      {
        path: '',
        name: 'Task',
        component: () => import('@/views/task/index'),
        meta: { title: '任务管理', icon: 'task' }
      }
    ]
  },
  // {
  //   path: '/test',
  //   component: Layout,
  //   children: [
  //     {
  //       path: '',
  //       name: 'Test',
  //       component: () => import('@/views/task/index'),
  //       meta: { title: '测试', icon: 'task' }
  //     }
  //   ]
  // },
  {
    path: '/hostgroup',
    component: Layout,
    children: [
      {
        path: '',
        name: 'HostGroup',
        component: () => import('@/views/hostgroup/index'),
        meta: { title: '主机组', icon: 'hostgroup' }
      }
    ]
  },
  {
    path: '/hosts',
    component: Layout,
    children: [
      {
        path: '',
        name: 'Host',
        component: () => import('@/views/hosts/index'),
        meta: { title: '主机', icon: 'host' }
      }
    ]
  },
  {
    path: '/log',
    component: Layout,
    children: [
      {
        path: '',
        name: 'Log',
        component: () => import('@/views/log/index'),
        meta: { title: '日志管理', icon: 'log' }
      }
    ]
  },
  {
    path: '/profile',
    component: Layout,
    children: [
      {
        path: '',
        name: 'Profile',
        component: () => import('@/views/profile/index'),
        meta: { title: '个人中心', icon: 'profile' }
      }
    ]
  },
  {
    path: '/notify',
    component: Layout,
    hidden: true,
    children: [
      {
        path: '',
        name: 'Notify',
        component: () => import('@/views/notify/index'),
        meta: { title: '通知消息', icon: 'profile' }
      }
    ]
  }

]

/**
 * asyncRoutes
 * the routes that need to be dynamically loaded based on user roles
 */
export const asyncRoutes = [
  {
    path: '/audit',
    component: Layout,
    children: [
      {
        path: '',
        name: 'Audit',
        component: () => import('@/views/audit/index'),
        meta: { title: '操作审计', icon: 'audit', roles: ['admin'] }
      }
    ]
  },
  {
    path: '/user',
    component: Layout,
    children: [
      {
        path: '',
        name: 'User',
        component: () => import('@/views/user/index'),
        meta: { title: '用户管理', icon: 'user', roles: ['admin'] }
      }
    ]
  },
  // 404 page must be placed at the end !!!
  { path: '*', redirect: '/404', hidden: true }
]

const createRouter = () => new Router({
  // mode: 'history', // require service support
  scrollBehavior: () => ({ y: 0 }),
  routes: constantRoutes
})

const router = createRouter()

// Detail see: https://github.com/vuejs/vue-router/issues/1234#issuecomment-357941465
export function resetRouter() {
  const newRouter = createRouter()
  router.matcher = newRouter.matcher // reset router
}

export default router
