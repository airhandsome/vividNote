<template>
  <div class="h-full flex items-center justify-center bg-gray-50">
    <div class="text-center">
      <h2 class="text-2xl font-bold text-gray-900 mb-4">
        欢迎使用 VividNote
      </h2>
      <p class="text-gray-500 mb-8">
        点击左侧「新建页面」开始写作，或选择已有页面继续编辑
      </p>
      <button
        @click="createNewPage"
        class="btn btn-primary inline-flex items-center gap-2"
      >
        <PlusIcon class="h-5 w-5" />
        新建页面
      </button>
    </div>
  </div>
</template>

<script setup>
import { PlusIcon } from '@heroicons/vue/24/outline'
import { useRouter } from 'vue-router'

const router = useRouter()

const createNewPage = async () => {
  try {
    const newPage = {
      title: '无标题',
      content: '',
      parentID: ''
    }
    const result = await window.go.main.App.CreatePage(newPage)
    if (result) {
      router.push(`/page/${result.id}`)
    }
  } catch (error) {
    console.error('创建页面失败:', error)
  }
}
</script> 