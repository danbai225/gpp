<template>
  <div class="container">
    <div class="center">
      <n-space vertical>
        <n-space>
          <n-progress type="circle" :percentage="percentageRef"/>
        </n-space>
        <n-space>
          <n-button style="margin-left: 20px" @click="start">
            开始加速
          </n-button>
          <n-button @click="stop">
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
const add = () => {
  percentageRef.value += 10
  if (percentageRef.value > 100) {
    percentageRef.value = 0
  }
}

const start = () => {
  Start().then(res => {
    console.log('res:', res)
  })
}
const stop = () => {
  Stop().then(res => {
    console.log('stopRes', res)
  })
}

const minus = () => {
  percentageRef.value -= 10
  if (percentageRef.value < 0) {
    percentageRef.value = 100
  }
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