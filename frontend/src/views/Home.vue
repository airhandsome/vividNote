<template>
  <div class="h-screen flex">
    <!-- 侧边栏 -->
    <div class="w-64 bg-white border-r border-gray-200 flex flex-col">
      <!-- 搜索框 -->
      <div class="p-4 border-b border-gray-200">
        <input
          type="text"
          placeholder="搜索..."
          class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>

      <!-- 笔记类型选择 -->
      <div class="flex flex-col">
        <button
          v-for="type in noteTypes"
          :key="type.id"
          @click="selectNoteType(type)"
          class="p-4 text-left hover:bg-gray-100 flex items-center justify-between"
          :class="{ 'bg-gray-100': currentType === type.id }"
        >
          <span>{{ type.name }}</span>
          <PlusIcon 
            @click.stop="createNewNote(type.id)"
            class="h-5 w-5 text-gray-500 hover:text-gray-700" 
          />
        </button>
      </div>

      <!-- 颜色选择器 (仅在原子笔记时显示) -->
      <div v-if="currentType === 'atomic'" class="p-4 flex gap-2">
        <div
          v-for="color in colors"
          :key="color"
          @click="selectColor(color)"
          class="w-6 h-6 rounded-full cursor-pointer"
          :class="{ 'ring-2 ring-offset-2 ring-blue-500': currentColor === color }"
          :style="{ backgroundColor: color }"
        ></div>
      </div>

      <!-- 笔记列表 -->
      <div class="flex-1 overflow-y-auto p-4">
        <TransitionGroup
          enter-active-class="transition duration-100 ease-out"
          enter-from-class="transform scale-95 opacity-0"
          enter-to-class="transform scale-100 opacity-100"
          leave-active-class="transition duration-75 ease-in"
          leave-from-class="transform scale-100 opacity-100"
          leave-to-class="transform scale-95 opacity-0"
        >
          <div
            v-for="page in filteredPages"
            :key="page.id"
            @click="openPage(page)"
            @dblclick="editPageTitle(page)"
            class="group p-3 rounded-lg cursor-pointer hover:bg-gray-100 transition-all duration-200 flex items-center gap-2"
            :class="{ 'bg-gray-100': currentPageId === page.id }"
          >
            <div 
              v-if="page.type === 'atomic'"
              class="w-3 h-3 rounded-full"
              :style="{ backgroundColor: page.color }"
            ></div>
            <span class="flex-1">{{ page.title || '无标题' }}</span>
            <button @click.stop="deletePage(page.id)" class="opacity-0 group-hover:opacity-100">
              <TrashIcon class="h-5 w-5 text-red-500 hover:text-red-700" />
            </button>
          </div>
        </TransitionGroup>
      </div>
    </div>

    <!-- 主内容区 -->
    <div class="flex-1 overflow-hidden bg-gray-50">
      <router-view 
        v-slot="{ Component }"
        :key="componentKey"
        @refresh-list="fetchPages"
      >
        <transition
          enter-active-class="transition-opacity duration-200"
          enter-from-class="opacity-0"
          enter-to-class="opacity-100"
          leave-active-class="transition-opacity duration-200"
          leave-from-class="opacity-100"
          leave-to-class="opacity-0"
        >
          <component 
            :is="Component" 
            @page-updated="handlePageUpdate"
            @refresh-list="fetchPages"
          />
        </transition>
      </router-view>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { PlusIcon, TrashIcon } from '@heroicons/vue/24/outline'

const route = useRoute()
const router = useRouter()
const pages = ref([])
const currentPageId = computed(() => route.params.id)

// 笔记类型
const noteTypes = [
  { id: 'original', name: '原始笔记' },
  { id: 'atomic', name: '原子笔记' }
]

// 原子笔记颜色
const colors = ['#FFB6C1', '#87CEEB', '#98FB98', '#DDA0DD', '#F0E68C']

const currentType = ref('original')
const currentColor = ref(colors[0])

// 根据当前类型过滤笔记
const filteredPages = computed(() => {
  return pages.value.filter(page => page.type === currentType.value)
})

// 选择笔记类型
const selectNoteType = (type) => {
  currentType.value = type.id
  if (type.id === 'atomic') {
    router.push('/atomic')
  }
  fetchPages()
}

// 选择颜色
const selectColor = (color) => {
  currentColor.value = color
}

const componentKey = ref(0)

// 创建新笔记
const createNewNote = async (type) => {
  try {
    const newPage = {
      title: '无标题',
      content: '',
      type: type,
      color: type === 'atomic' ? currentColor.value : null
    }
    const result = await window.go.main.App.CreatePage(newPage)
    if (result) {
      await fetchPages()
      if (type === 'atomic') {
        // 强制刷新 AtomicNoteBoard 组件
        componentKey.value += 1
        router.push('/atomic')
      } else {
        router.push(`/original/${result.id}`)
      }
    }
  } catch (error) {
    console.error('创建页面失败:', error)
  }
}

// 打开页面
const openPage = (page) => {
  if (page.type === 'atomic') {
    router.push({ 
      path: '/atomic',
      query: { selected: page.id }  // 添加查询参数
    })
  } else {
    router.push(`/original/${page.id}`)
  }
}

// 获取页面列表
const fetchPages = async () => {
  try {
    console.log('currentType.value', currentType.value)
    const result = await window.go.main.App.GetChildPages(currentType.value)
    pages.value = result
  } catch (error) {
    console.error('获取页面列表失败:', error)
  }
}

// 处理页面更新
const handlePageUpdate = (updatedPage) => {
  const index = pages.value.findIndex(p => p.id === updatedPage.id)
  if (index !== -1) {
    pages.value[index] = { ...pages.value[index], ...updatedPage }
    // 强制更新数组以触发视图更新
    pages.value = [...pages.value]
  }
}

// 删除页面
const deletePage = async (pageId) => {
  try {
    if (!pageId) {
      console.error('Page ID is required for deletion')
      return
    }
    await window.go.main.App.DeletePage(pageId)
    await fetchPages()
    if (currentPageId.value === pageId) {
      router.push('/')
    }
  } catch (error) {
    console.error('删除页面失败:', error)
  }
}

// 编辑页面标题
const editPageTitle = async (page) => {
  const newTitle = prompt('请输入新的标题', page.title)
  if (newTitle !== null && newTitle !== page.title) {
    page.title = newTitle
    await window.go.main.App.UpdatePage(page)
    handlePageUpdate(page)
    if (currentPageId.value === page.id) {
      // 更新右侧编辑器的标题
      const editorComponent = router.currentRoute.value.matched[0].instances.default
      if (editorComponent) {
        editorComponent.title = newTitle
      }
    }
  }
}

onMounted(() => {
  fetchPages()
})
</script>

<style scoped>
/* 样式示例 */
.page-atom {
  background-color: #f0f0f0;
  border: 1px solid #ccc;
}

.page-original {
  background-color: #fff;
  border: 1px solid #ddd;
}
</style>