<template>
  <div class="min-h-screen p-8">
    <!-- 搜索框 -->
    <div class="mb-6">
      <input
        v-model="searchQuery"
        type="text"
        placeholder="搜索..."
        class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
      />
    </div>

    <!-- 原子笔记展示区域 -->
    <div class="grid grid-cols-3 gap-6">
      <div
        v-for="note in filteredNotes"
        :key="note.id"
        :id="`note-${note.id}`"
        class="group relative p-6 rounded-lg shadow-md transition-all duration-300"
        :class="{
          'ring-4 ring-blue-500 scale-105': selectedNoteId === note.id,
          'hover:shadow-lg': selectedNoteId !== note.id
        }"
        :style="{ 
          backgroundColor: note.color,
          transform: selectedNoteId === note.id ? 'scale(1.05)' : 'scale(1)'
        }"
      >
        <!-- 操作按钮 -->
        <div class="absolute top-2 right-2 flex gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
          <button 
            @click.stop="changeColor(note)"
            class="p-2 bg-white rounded-full shadow hover:shadow-md"
            title="修改颜色"
          >
            <SwatchIcon class="h-4 w-4 text-gray-600" />
          </button>
          <button 
            @click.stop="deleteNote(note.id)"
            class="p-2 bg-white rounded-full shadow hover:shadow-md"
            title="删除笔记"
          >
            <TrashIcon class="h-4 w-4 text-red-500" />
          </button>
        </div>

        <!-- 笔记内容 -->
        <div @click="selectNote(note.id)" class="cursor-pointer">
          <h3 class="text-lg font-semibold mb-2">{{ note.title || '无标题' }}</h3>
          <div class="prose prose-sm" v-html="note.content"></div>
        </div>
      </div>
    </div>

    <!-- 编辑弹窗 -->
    <div
      v-if="showEditModal"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
      @click.self="closeEditModal"
    >
      <div 
        class="bg-white rounded-lg p-6 w-full max-w-lg"
        :style="{ backgroundColor: currentNote?.color }"
      >
        <input
          v-model="currentNote.title"
          class="w-full mb-4 px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white"
          placeholder="标题"
        />
        <textarea
          v-model="currentNote.content"
          class="w-full h-48 px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white"
          placeholder="内容"
        ></textarea>
        
        <!-- 颜色选择器 -->
        <div class="flex gap-2 mt-4 mb-4">
          <div
            v-for="color in colors"
            :key="color"
            @click="currentNote.color = color"
            class="w-6 h-6 rounded-full cursor-pointer border-2"
            :class="{ 'border-blue-500': currentNote.color === color, 'border-transparent': currentNote.color !== color }"
            :style="{ backgroundColor: color }"
          ></div>
        </div>

        <div class="flex justify-end gap-4 mt-4">
          <button
            @click="closeEditModal"
            class="px-4 py-2 text-gray-600 hover:text-gray-800"
          >
            取消
          </button>
          <button
            @click="saveNote"
            class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          >
            保存
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { SwatchIcon, TrashIcon } from '@heroicons/vue/24/outline'

const route = useRoute()
const atomicNotes = ref([])
const showEditModal = ref(false)
const currentNote = ref(null)
const searchQuery = ref('')
const selectedNoteId = ref(null)
const emit = defineEmits(['refresh-list'])
// 监听路由查询参数变化
watch(() => route.query.selected, (newId) => {
  if (newId) {
    selectedNoteId.value = newId
    // 滚动到选中的笔记
    scrollToNote(newId)
  }
}, { immediate: true })

// 滚动到选中的笔记
const scrollToNote = (noteId) => {
  nextTick(() => {
    const element = document.getElementById(`note-${noteId}`)
    if (element) {
      element.scrollIntoView({ 
        behavior: 'smooth', 
        block: 'center' 
      })
    }
  })
}

// 选择笔记
const selectNote = (noteId) => {
  selectedNoteId.value = noteId
  editNote(atomicNotes.value.find(note => note.id === noteId))
}

// 预定义的颜色选项
const colors = [
  '#FFB6C1', // 浅粉红
  '#87CEEB', // 天蓝色
  '#98FB98', // 浅绿色
  '#DDA0DD', // 梅红色
  '#F0E68C', // 卡其色
  '#E6E6FA', // 淡紫色
  '#FFA07A', // 浅鲑鱼色
  '#B0E0E6'  // 粉蓝色
]

// 获取所有原子笔记
const fetchAtomicNotes = async () => {
  try {
    const result = await window.go.main.App.GetChildPages( 'atomic')
    atomicNotes.value = result
    // 如果有选中的笔记ID，检查它是否还存在
    if (selectedNoteId.value && !result.some(note => note.id === selectedNoteId.value)) {
      selectedNoteId.value = null
    }
  } catch (error) {
    console.error('获取原子笔记失败:', error)
  }
}

// 过滤笔记
const filteredNotes = computed(() => {
  return atomicNotes.value.filter(note =>
    note.title.includes(searchQuery.value) || note.content.includes(searchQuery.value)
  )
})

// 编辑笔记
const editNote = (note) => {
  currentNote.value = { ...note }
  showEditModal.value = true
}

// 关闭编辑弹窗
const closeEditModal = () => {
  showEditModal.value = false
  currentNote.value = null
}

// 保存笔记
const saveNote = async () => {
  try {
    await window.go.main.App.UpdatePage(currentNote.value)
    const index = atomicNotes.value.findIndex(note => note.id === currentNote.value.id)
    if (index !== -1) {
      atomicNotes.value[index] = { ...currentNote.value }
    }
    emit('refresh-list')
    closeEditModal()
  } catch (error) {
    console.error('保存笔记失败:', error)
  }
}

// 删除笔记
const deleteNote = async (noteId) => {
  try {
    await window.go.main.App.DeletePage(null, noteId)
    atomicNotes.value = atomicNotes.value.filter(note => note.id !== noteId)
    // 触发刷新事件
    emit('refresh-list')
    // 重新获取笔记列表
    fetchAtomicNotes()
  } catch (error) {
    console.error('删除笔记失败:', error)
  }
}

// 修改改变颜色的方法
const changeColor = (note) => {
  currentNote.value = { ...note }
  showEditModal.value = true
}

onMounted(() => {
  fetchAtomicNotes()
})
</script>

<style scoped>
.grid {
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
}

/* 添加平滑过渡效果 */
.group {
  transform-origin: center;
  transition: all 0.3s ease-in-out;
}

/* 选中效果的动画 */
@keyframes highlight {
  0% {
    transform: scale(1);
    box-shadow: 0 0 0 0 rgba(59, 130, 246, 0.5);
  }
  50% {
    transform: scale(1.5);
    box-shadow: 0 0 0 10px rgba(59, 130, 246, 0);
  }
  100% {
    transform: scale(1.1);
    box-shadow: 0 0 0 0px rgba(59, 130, 246, 0);
    border: 5px solid rgba(253, 3, 103, 0.5);
  }
}

.scale-105 {
  animation: highlight 0.5s ease-out forwards;
}
</style> 