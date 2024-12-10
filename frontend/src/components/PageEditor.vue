<template>
  <div class="min-h-screen" :style="{ backgroundColor }">
    <!-- 工具栏 -->
    <div class="sticky top-0 z-10 bg-white/80 backdrop-blur-sm border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="py-4 flex items-center justify-between">
          <div class="flex items-center gap-4">
            <button
              @click="savePage"
              class="btn btn-primary flex items-center gap-2"
              :disabled="isSaving"
            >
              <DocumentCheckIcon v-if="!isSaving" class="h-5 w-5" />
              <ArrowPathIcon v-else class="h-5 w-5 animate-spin" />
              {{ isSaving ? '保存中...' : '保存' }}
            </button>
            
            <Menu as="div" class="relative">
              <MenuButton class="btn btn-secondary flex items-center gap-2">
                <SwatchIcon class="h-5 w-5" />
                背景色
              </MenuButton>
              <MenuItems class="absolute left-0 mt-2 w-56 bg-white rounded-lg shadow-lg p-2">
                <input
                  type="color"
                  v-model="backgroundColor"
                  @change="saveBackgroundColor"
                  class="w-full h-8 rounded"
                />
              </MenuItems>
            </Menu>

            <button
              @click="exportToPDF"
              class="btn btn-secondary flex items-center gap-2"
            >
              <ArrowDownTrayIcon class="h-5 w-5" />
              导出 PDF
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 编辑区域 -->
    <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <input
        v-model="title"
        class="w-full text-4xl font-bold bg-transparent border-none outline-none mb-8 placeholder-gray-400"
        placeholder="无标题"
      />
      
      <div class="prose prose-lg max-w-none bg-white rounded-lg shadow-sm">
        <tiptap-editor
          v-model="content"
          @update:modelValue="handleContentUpdate"
        />
      </div>
    </div>

    <!-- PDF 导出区域 (隐藏) -->
    <div ref="pdfContent" class="pdf-content">
      <div class="p-8 bg-white">
        <h1 class="text-4xl font-bold mb-8">{{ title || '无标题' }}</h1>
        <div 
          class="prose prose-lg max-w-none"
          v-html="sanitizedContent"
        ></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, onBeforeUnmount, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Menu, MenuButton, MenuItems } from '@headlessui/vue'
import { 
  DocumentCheckIcon, 
  ArrowPathIcon, 
  SwatchIcon,
  ArrowDownTrayIcon 
} from '@heroicons/vue/24/outline'
import TiptapEditor from './Editor/TiptapEditor.vue'
import html2pdf from 'html2pdf.js'
import DOMPurify from 'dompurify'

const route = useRoute()
const router = useRouter()
const pdfContent = ref(null)
const title = ref('')
const content = ref('')
const currentPageId = ref('')
const backgroundColor = ref('#ffffff')
const isSaving = ref(false)
const hasUnsavedChanges = ref(false)

const emit = defineEmits(['pageUpdated'])

// 加载页面内容
const loadPage = async (id) => {
  if (!id) return
  
  try {
    const page = await window.go.main.App.GetPage(id)
    if (page) {
      title.value = page.title || ''
      content.value = page.content || ''
      backgroundColor.value = page.background_color || '#ffffff'
      currentPageId.value = id
      hasUnsavedChanges.value = false
    }
  } catch (error) {
    console.error('加载页面失败:', error)
  }
}

// 保存页面
const savePage = async () => {
  if (!currentPageId.value || isSaving.value) return
  
  try {
    isSaving.value = true
    const pageData = {
      id: currentPageId.value,
      title: title.value,
      content: content.value,
      background_color: backgroundColor.value,
      type: 'original',      
    }
    await window.go.main.App.UpdatePage(pageData)
    emit('pageUpdated', pageData)
    hasUnsavedChanges.value = false
  } catch (error) {
    console.error('保存失败:', error)
  } finally {
    isSaving.value = false
  }
}

// 清理 HTML 内容
const sanitizedContent = computed(() => {
  return DOMPurify.sanitize(content.value, {
    ADD_TAGS: ['h1', 'h2', 'h3', 'p', 'a', 'ul', 'ol', 'li', 'blockquote', 'code', 'pre'],
    ADD_ATTR: ['href', 'target', 'rel', 'style']
  })
})

// 修改导出 PDF 函数
const exportToPDF = async () => {
  try {
    // 先保存当前内容
    if (hasUnsavedChanges.value) {
      await savePage()
    }

    if (!pdfContent.value) {
      console.error('PDF 内容区域未找到')
      return
    }

    // 等待内容渲染
    await new Promise(resolve => setTimeout(resolve, 500))

    const element = pdfContent.value
    const opt = {
      margin: [15, 15],
      filename: `${title.value || '无标题'}.pdf`,
      image: { type: 'jpeg', quality: 0.98 },
      html2canvas: { 
        scale: 2,
        useCORS: true,
        logging: true,
        width: 794, // A4 宽度 (px)
        windowWidth: 794,
      },
      jsPDF: { 
        unit: 'mm', 
        format: 'a4', 
        orientation: 'portrait'
      }
    }

    await html2pdf().set(opt).from(element).save()
  } catch (error) {
    console.error('导出 PDF 失败:', error)
  }
}

// 内容更新时标记未保存
const handleContentUpdate = (newContent) => {
  content.value = newContent
  hasUnsavedChanges.value = true
}

// 监听标题变化
watch(title, () => {
  hasUnsavedChanges.value = true
})

// 监听路由变化，自动保存
watch(
  () => route.params.id,
  async (newId, oldId) => {
    if (oldId && hasUnsavedChanges.value) {
      await savePage()
    }
    if (newId) {
      loadPage(newId)
    }
  },
  { immediate: true }
)

// 页面关闭前保存
onBeforeUnmount(async () => {
  if (hasUnsavedChanges.value) {
    await savePage()
  }
})
</script>

<style>
.pdf-content {
  position: fixed;
  left: -9999px;
  top: -9999px;
  width: 794px; /* A4 宽度 */
  background: white;
  z-index: -1;
}

/* PDF 导出样式 */
.pdf-content h1 {
  font-size: 28px;
  margin-bottom: 20px;
  color: #000;
}

.pdf-content .prose {
  max-width: none;
  color: #000;
}

.pdf-content .prose h1 { font-size: 24px; margin: 24px 0; }
.pdf-content .prose h2 { font-size: 20px; margin: 20px 0; }
.pdf-content .prose h3 { font-size: 18px; margin: 16px 0; }
.pdf-content .prose p { margin: 16px 0; }
.pdf-content .prose ul,
.pdf-content .prose ol { margin: 16px 0; padding-left: 24px; }
.pdf-content .prose blockquote { 
  border-left: 4px solid #e5e7eb;
  padding-left: 16px;
  margin: 16px 0;
}

.page-container {
  min-height: 100vh;
  transition: background-color 0.3s;
}

.editor-wrapper {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;
}

.toolbar {
  position: sticky;
  top: 0;
  padding: 1rem;
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(10px);
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #eaeaea;
  z-index: 100;
}

.save-btn {
  padding: 0.5rem 1.5rem;
  background-color: #2ecc71;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 500;
  transition: background-color 0.2s;
}

.save-btn:hover {
  background-color: #27ae60;
}

.color-picker {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.color-picker label {
  font-size: 0.9rem;
  color: #666;
}

.color-picker input[type="color"] {
  padding: 0;
  width: 40px;
  height: 30px;
  border: 1px solid #ddd;
  border-radius: 4px;
  cursor: pointer;
}

.title-input {
  font-size: 2rem;
  font-weight: bold;
  width: 100%;
  border: none;
  outline: none;
  margin-bottom: 1.5rem;
  background: transparent;
  padding: 0.5rem;
  border-radius: 4px;
}

.title-input:focus {
  background: rgba(255, 255, 255, 0.8);
}

.title-input::placeholder {
  color: #adb5bd;
}

.pdf-content {
  display: none;
}

@media print {
  .pdf-content {
    display: block;
  }
}
</style> 