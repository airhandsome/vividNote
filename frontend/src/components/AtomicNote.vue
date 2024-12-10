<template>
  <div class="min-h-screen p-8" :style="{ backgroundColor: page.color || '#ffffff' }">
    <div class="bg-white p-6 rounded-lg shadow-md">
      <h2 class="text-xl font-bold mb-4">{{ page.title }}</h2>
      <div class="prose max-w-none" v-html="page.content"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const page = ref({
  title: '',
  content: '',
  color: ''
})

// 获取页面数据
const fetchPage = async () => {
  try {
    const result = await window.go.main.App.GetPage(route.params.id)
    if (result) {
      page.value = result
    }
  } catch (error) {
    console.error('获取页面失败:', error)
  }
}

onMounted(() => {
  fetchPage()
})
</script> 