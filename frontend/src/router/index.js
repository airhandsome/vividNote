import { createRouter, createWebHashHistory } from 'vue-router'
import Home from '../views/Home.vue'
import PageEditor from '../components/PageEditor.vue'
import AtomicNoteBoard from '../components/AtomicNoteBoard.vue'

const routes = [
  {
    path: '/',
    component: Home,
    children: [
      {
        path: 'original/:id',
        name: 'page',
        component: PageEditor,
        props: true
      },
      {
        path: 'atomic',
        name: 'atomic',
        component: AtomicNoteBoard
      },
      {
        path: '',
        name: 'welcome',
        component: () => import('../views/Welcome.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router 