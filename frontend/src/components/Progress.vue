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
          <n-button :disabled="btnDisabled" @click="!state?start():stop()" style="margin-left: 20px">
            {{ btnText }}
          </n-button>
        </n-space>
      </n-space>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {ref} from 'vue'
import {Start, Stop} from "../../wailsjs/go/main/App";

const percentageRef = ref(0)
const state = ref(false)
const btnText = ref('开始加速')
const btnDisabled = ref(false)

const start = () => {
  btnDisabled.value = true
  btnText.value = '加速中.'
  Start().then(res => {
    state.value = true
    console.log('startRes', res)
    let timer = setInterval(() => {
      percentageRef.value += 10
      if (percentageRef.value === 100) {
        clearInterval(timer)
        btnText.value = '结束加速'
        btnDisabled.value = false
      }
    }, 100)
  })
}
const stop = () => {
  Stop().then(res => {
    percentageRef.value = 0
    state.value = false
    btnText.value = '开始加速'
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