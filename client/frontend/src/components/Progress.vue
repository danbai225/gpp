<template>
  <div class="container">
    <div class="center">
      <n-space vertical>
        <n-space>
          <n-progress type="circle" :percentage="percentage"/>
        </n-space>
        <n-space>
          <n-button style="margin-left: 16px">
            开始加速
          </n-button>
        </n-space>
      </n-space>
    </div>

  </div>
</template>

<script lang="ts">
import {defineComponent, ref} from 'vue'

export default defineComponent({
  setup() {
    const percentageRef = ref(0)
    const add = () => {
      percentageRef.value += 10
      if (percentageRef.value > 100) {
        percentageRef.value = 0
      }
    }

    const minus = () => {
      percentageRef.value -= 10
      if (percentageRef.value < 0) {
        percentageRef.value = 100
      }
    }

    return {
      percentage: percentageRef,
      add,
      minus
    }
  }
})
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