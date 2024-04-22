<template>
  <div class="container">
    <div class="center">
      <n-space vertical>
        <n-space>
          <n-progress type="circle"
                      :status="percentageRef<=25?'error':percentageRef<=50?'warning':percentageRef<=75?'info':'success'"
                      :percentage="percentageRef">
            {{ percentageRef === 100 ? '加速完成' : percentageRef === 0 ? '未开始' : '正在加速' }}
          </n-progress>
        </n-space>
        <n-space>
          <n-button v-if="!state" @click="start" style="margin-left: 20px">
            开始加速
          </n-button>
          <n-button v-else @click="stop" style="margin-left: 20px">
            结束加速
          </n-button>
        </n-space>
      </n-space>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {defineComponent, ref} from 'vue'
import {Start, Stop} from "../../wailsjs/go/main/App";

const percentageRef = ref(0)
const state = ref(false)

const start = () => {
  setTimeout(() => {
    let i = 0
    const timer = setInterval(() => {
      i += 20
      percentageRef.value = i
      if (i >= 100) {
        clearInterval(timer)
      }
    }, 500)
  }, 1000)
  Start().then(res => {
    console.log('startRes', res)
    if (res === 'ok') {
      state.value = true
    }
  })
}
const stop = () => {
  Stop().then(res => {
    percentageRef.value = 0
    state.value = false
    console.log('stopRes', res)
  })
}


</script>

<style>
.container {
  display: flex;
  justify-content: center;
  align-items: center;
}

.center {
  /* 可以添加宽度、高度等样式 */
  margin-top: 20%;
}
</style>