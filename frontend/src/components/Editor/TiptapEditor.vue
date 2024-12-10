<template>
  <div class="editor-wrapper">
    <div class="editor-toolbar">
      <div class="toolbar-group">
        <select 
          class="toolbar-select"
          @change="setFontSize"
          :value="currentFontSize"
          title="字号"
        >
          <option value="small">小号</option>
          <option value="normal">正常</option>
          <option value="large">大号</option>
          <option value="huge">超大</option>
        </select>
      </div>

      <div class="toolbar-divider" />
      
      <div class="toolbar-group">
        <button
          v-for="heading in [1, 2, 3]"
          :key="heading"
          @click="editor?.chain().focus().toggleHeading({ level: heading }).run()"
          :class="{ 'is-active': editor?.isActive('heading', { level: heading }) }"
          class="toolbar-btn"
          :title="`标题 ${heading} (大标题)`"
        >
          H{{ heading }}
        </button>
        <button
          @click="editor?.chain().focus().setParagraph().run()"
          :class="{ 'is-active': editor?.isActive('paragraph') }"
          class="toolbar-btn"
          title="正文段落"
        >
          <DocumentTextIcon class="w-4 h-4" />
        </button>
      </div>

      <div class="toolbar-divider" />

      <div class="toolbar-group">
        <button
          @click="editor?.chain().focus().toggleBold().run()"
          :class="{ 'is-active': editor?.isActive('bold') }"
          class="toolbar-btn"
          title="加粗 (Ctrl+B)"
        >
          <BoltIcon class="w-4 h-4" />
        </button>
        <button
          @click="editor?.chain().focus().toggleItalic().run()"
          :class="{ 'is-active': editor?.isActive('italic') }"
          class="toolbar-btn"
          title="斜体 (Ctrl+I)"
        >
          <ChartBarIcon class="w-4 h-4" />
        </button>
        <button
          @click="editor?.chain().focus().toggleStrike().run()"
          :class="{ 'is-active': editor?.isActive('strike') }"
          class="toolbar-btn"
          title="删除线"
        >
          <MinusIcon class="w-4 h-4" />
        </button>
        <button
          @click="editor?.chain().focus().toggleCode().run()"
          :class="{ 'is-active': editor?.isActive('code') }"
          class="toolbar-btn"
          title="行内代码"
        >
          <CodeBracketIcon class="w-4 h-4" />
        </button>
      </div>

      <div class="toolbar-divider" />

      <div class="toolbar-group">
        <button
          @click="editor?.chain().focus().toggleBulletList().run()"
          :class="{ 'is-active': editor?.isActive('bulletList') }"
          class="toolbar-btn"
          title="无序列表"
        >
          <ListBulletIcon class="w-4 h-4" />
        </button>
        <button
          @click="editor?.chain().focus().toggleOrderedList().run()"
          :class="{ 'is-active': editor?.isActive('orderedList') }"
          class="toolbar-btn"
          title="有序列表"
        >
          <QueueListIcon class="w-4 h-4" />
        </button>
        <button
          @click="editor?.chain().focus().toggleTaskList().run()"
          :class="{ 'is-active': editor?.isActive('taskList') }"
          class="toolbar-btn"
          title="任务列表"
        >
          <CheckCircleIcon class="w-4 h-4" />
        </button>
      </div>

      <div class="toolbar-divider" />

      <div class="toolbar-group">
        <button
          @click="editor?.chain().focus().toggleBlockquote().run()"
          :class="{ 'is-active': editor?.isActive('blockquote') }"
          class="toolbar-btn"
          title="引用"
        >
          <ChatBubbleLeftIcon class="w-4 h-4" />
        </button>
        <button
          @click="editor?.chain().focus().toggleCodeBlock().run()"
          :class="{ 'is-active': editor?.isActive('codeBlock') }"
          class="toolbar-btn"
          title="代码块"
        >
          <CodeBracketSquareIcon class="w-4 h-4" />
        </button>
        <button
          @click="addImage"
          class="toolbar-btn"
          title="插入图片"
        >
          <PhotoIcon class="w-4 h-4" />
        </button>
      </div>

      <div class="toolbar-divider" />

      <div class="toolbar-group">
        <button
          @click="editor?.chain().focus().undo().run()"
          :disabled="!editor?.can().undo()"
          class="toolbar-btn"
          title="撤销 (Ctrl+Z)"
        >
          <ArrowUturnLeftIcon class="w-4 h-4" />
        </button>
        <button
          @click="editor?.chain().focus().redo().run()"
          :disabled="!editor?.can().redo()"
          class="toolbar-btn"
          title="重做 (Ctrl+Y)"
        >
          <ArrowUturnRightIcon class="w-4 h-4" />
        </button>
      </div>
    </div>
    
    <editor-content 
      v-if="editor" 
      :editor="editor" 
      class="editor-content"
    />
  </div>
</template>

<script setup>
import { ref, onBeforeUnmount } from 'vue'
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import TextStyle from '@tiptap/extension-text-style'
import { Extension } from '@tiptap/core'
import {
  DocumentTextIcon,
  ListBulletIcon,
  QueueListIcon,
  CheckCircleIcon,
  PhotoIcon,
  ArrowUturnLeftIcon,
  ArrowUturnRightIcon,
  CodeBracketIcon,
  CodeBracketSquareIcon,
  ChatBubbleLeftIcon,
  BoltIcon,
  ChartBarIcon,
  MinusIcon
} from '@heroicons/vue/24/outline'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue'])

const currentFontSize = ref('normal')

// 创建字号扩展
const CustomFontSize = Extension.create({
  name: 'fontSize',
  addAttributes() {
    return {
      fontSize: {
        default: 'normal',
        parseHTML: element => element.style.fontSize,
        renderHTML: attributes => {
          return { style: `font-size: ${attributes.fontSize || '1em'}` }
        }
      }
    }
  },
  addCommands() {
    return {
      setFontSize: (size) => ({ chain }) => {
        const headingLevels = {
          huge: 1,
          large: 2,
          normal: 3,
          small: null // 正文
        }
        const level = headingLevels[size]
        if (level) {
          return chain().toggleHeading({ level }).run()
        } else {
          return chain().setParagraph().run()
        }
      }
    }
  }
})

const editor = useEditor({
  content: props.modelValue,
  extensions: [
    StarterKit,
    TextStyle,
    CustomFontSize
  ],
  editorProps: {
    attributes: {
      class: 'prose prose-lg focus:outline-none max-w-none'
    }
  },
  onUpdate: ({ editor }) => {
    emit('update:modelValue', editor.getHTML())
  }
})

const addImage = () => {
  const url = window.prompt('输入图片URL')
  if (url) {
    editor.value?.chain().focus().setImage({ src: url }).run()
  }
}

const setFontSize = (event) => {
  const size = event.target.value
  currentFontSize.value = size
  editor.value?.chain().focus().setFontSize(size).run()
}

// 确保在组件卸载时销毁编辑器实例
onBeforeUnmount(() => {
  if (editor.value) {
    editor.value.destroy()
  }
})
</script>

<style>
.editor-wrapper {
  @apply border border-gray-200 rounded-lg overflow-hidden bg-white;
}

.editor-toolbar {
  @apply flex items-center gap-2 p-2 border-b border-gray-200 bg-gray-50 flex-wrap;
}

.toolbar-group {
  @apply flex items-center gap-1;
}

.toolbar-divider {
  @apply w-px h-6 bg-gray-200 mx-2;
}

.toolbar-btn {
  @apply p-1.5 rounded text-sm text-gray-600 hover:bg-gray-200 transition-colors relative;
}

.toolbar-btn.is-active {
  @apply bg-gray-200 text-gray-900;
}

.toolbar-btn:disabled {
  @apply opacity-50 cursor-not-allowed;
}

.editor-content {
  @apply p-4 min-h-[200px];
}

.ProseMirror {
  @apply outline-none;
}

.ProseMirror p.is-editor-empty:first-child::before {
  content: attr(data-placeholder);
  @apply text-gray-400 float-left h-0 pointer-events-none;
}

/* 任务列表样式 */
ul[data-type="taskList"] {
  @apply list-none p-0;
}

ul[data-type="taskList"] li {
  @apply flex gap-2 items-start;
}

ul[data-type="taskList"] li > label {
  @apply mt-1;
}

/* 引用样式 */
blockquote {
  @apply border-l-4 border-gray-300 pl-4 italic;
}

/* 代码块样式 */
pre {
  @apply bg-gray-900 text-white p-4 rounded-lg overflow-x-auto;
}

code {
  @apply bg-gray-100 px-1.5 py-0.5 rounded text-sm;
}

.toolbar-btn:hover::after {
  content: attr(title);
  @apply absolute left-1/2 -translate-x-1/2 bottom-full mb-2 px-2 py-1 
    text-xs text-white bg-gray-800 rounded whitespace-nowrap z-50;
}

.toolbar-btn:hover::before {
  content: '';
  @apply absolute left-1/2 -translate-x-1/2 bottom-full mb-1
    border-4 border-transparent border-t-gray-800;
}

/* 添加字号选择器样式 */
.toolbar-select {
  @apply px-2 py-1.5 rounded text-sm text-gray-600 bg-white border border-gray-200 
    hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500;
  min-width: 80px;
}

.toolbar-select option {
  @apply py-1;
}

/* 字号样式 */
.ProseMirror {
  .text-small {
    font-size: 0.875em;
  }
  
  .text-normal {
    font-size: 1em;
  }
  
  .text-large {
    font-size: 1.25em;
  }
  
  .text-huge {
    font-size: 1.5em;
  }
}

/* 添加标题样式 */
.ProseMirror h1 {
  font-size: 2.5em !important;
  margin-bottom: 1em;
}

.ProseMirror h2 {
  font-size: 2em !important;
  margin-bottom: 0.8em;
}

.ProseMirror h3 {
  font-size: 1.5em !important;
  margin-bottom: 0.6em;
}
</style> 