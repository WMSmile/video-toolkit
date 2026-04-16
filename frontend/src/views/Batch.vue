<template>
  <div>
    <h2>批量视频转MP3</h2>
    <button @click="start" class="p-2 bg-purple-600 text-white">选择文件夹并转换</button>
    <div class="mt-4 p-3 bg-gray-100 h-48 overflow-y-auto">{{ log }}</div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
const log = ref('')
import {Events} from "@wailsio/runtime";
import {ConvertService, FileService} from "../../bindings/video-toolkit/services";

async function start() {
  const path = await FileService.SelectDirectory('选择要转换的文件夹')
  const res = await ConvertService.ConvertFolder(path)
  log.value = res.join('\n')
}
</script>
