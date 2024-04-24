<template>
  <div class="container">
    <div class="center">
      <n-space vertical>
        <n-space style="text-align: center">

          <n-progress
              type="circle"
              :height="24"
              :status="percentageRef<=25?'error':percentageRef<=50?'warning':percentageRef<=75?'info':'success'"
              :percentage="percentageRef"
          >
            <n-space vertical size="small">
              {{ percentageRef === 100 ? '加速完成' : percentageRef === 0 ? '未开始' : '正在加速' }}
              <n-space vertical size="small" v-if="showGameHttpInfo">
                <p @click="getList()">
                  Game:{{ gamePeer === undefined ? '未选择' : gamePeer.name }}
                  <n-gradient-text v-if="gamePeer"
                                   :type="gamePeer.ping<60?'success':gamePeer.ping<100?'warning':'error'">
                    {{ gamePeer.ping }}
                  </n-gradient-text>
                </p>
                <p @click="getList()">
                  Http:{{ httpPeer === undefined ? '未选择' : httpPeer.name }}
                  <n-gradient-text v-if="httpPeer"
                                   :type="httpPeer.ping<60?'success':httpPeer.ping<100?'warning':'error'">
                    {{ httpPeer.ping }}
                  </n-gradient-text>
                </p>
              </n-space>
              <n-space vertical size="small" v-if="showUpDowInfo">
                <p>
                  上传:
                  <n-gradient-text v-if="up" type="success">
                    {{ (up).toFixed(2) }}kb
                  </n-gradient-text>
                </p>
                <p>
                  下载:
                  <n-gradient-text v-if="down" type="success">
                    {{ (down).toFixed(2) }}kb
                  </n-gradient-text>
                </p>
              </n-space>
            </n-space>
          </n-progress>
        </n-space>
        <n-space>
          <n-button :disabled="btnDisabled" @click="!state?start():stop()" style="margin-left: 55px">
            {{ btnText }}
          </n-button>
          <!--          <n-button @click="getList()">-->
          <!--            list-->
          <!--          </n-button>-->
          <!--          <n-button @click="getStatus()">-->
          <!--            Status-->
          <!--          </n-button>-->
        </n-space>
      </n-space>
      <div>
        <n-modal
            v-model:show="showModal"
            :mask-closable="false"
            preset="dialog"
            title="节点列表"
            positive-text="确认"
            negative-text="取消"
            @positive-click="submitCallback"
        >
          <n-select
              v-model:value="gameValue"
              vertical
              filterable
              :options="gameOpt"
              placeholder="请选择Game"
              value-field="val"
              label-field="name"
          />
          <br>
          <n-select
              v-model:value="httpValue"
              vertical
              filterable
              :options="gameOpt"
              placeholder="请选择Http"
              value-field="val"
              label-field="name"
          />
          <br>
          <n-input
              v-model:value="newUrl"
              type="textarea"
              placeholder="导入新连接"
          />
        </n-modal>

      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {ref, defineComponent, Ref, reactive, onMounted} from 'vue'
import {Add, List, SetPeer, Start, Status, Stop} from "../../wailsjs/go/main/App";
import {SelectOption, SelectGroupOption} from 'naive-ui'
import {config, data} from "../../wailsjs/go/models";
import {onBeforeMount} from "@vue/runtime-core";

const percentageRef = ref(0)
const state = ref(false)
const btnText = ref('开始加速')
const btnDisabled = ref(false)
const showModal = ref(false)
const httpDialog = ref(false)
const gameOpt = ref(Array<SelectOption | SelectGroupOption>())
const httpOpt = ref(Array<SelectOption | SelectGroupOption>())
const gameValue = ref()
const httpValue = ref()

const gamePeer: Ref<any> | undefined = ref()
const httpPeer: Ref<any> | undefined = ref()
const up = ref()
const down = ref()

const showGameHttpInfo = ref(true)
const showUpDowInfo = ref(false)

const newUrl = ref()

let time = ref()
onMounted(() => {
  getStatus()
  time.value = setInterval(() => {
    getStatus()
  }, 1000);
})

onBeforeMount(() => {
  clearInterval(time.value)
  time.value = null;
})

const getStarInfo = () => {
  Start().then(res => {
    console.log(res)
    if (res === "ok") {
      showGameHttpInfo.value = false
    } else if (res === "running") {
      showGameHttpInfo.value = true
      return
    }
  })
}

const start = () => {
  btnDisabled.value = true
  showGameHttpInfo.value = false
  showUpDowInfo.value = true
  btnText.value = '加速中.'
  Start().then(res => {
    state.value = true
    // console.log('startRes', res)
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
    showGameHttpInfo.value = true
    showUpDowInfo.value = false
    btnText.value = '开始加速'
    // console.log('stopRes', res)
  })
}
// Z3BwOi8vdmxlc3NAMTIzLjU4LjIxMi4xOTU6MzQ1NTYvYmFkYjE3ZWYtZWIyMi00ZTAzLTliMTctZWZlYjIyNGUwM2U3
const getList = () => {
  showModal.value = true
  gameOpt.value = Array<SelectOption | SelectGroupOption>()
  List().then(res => {
    res.forEach((item) => {
      gameOpt.value.push({
        name: item.name + '-' + item.ping + 'ms',
        val: item.name
      })
    })
  })
}

const getStatus = () => {
  Status().then(res => {
    gamePeer.value = res.game_peer
    httpPeer.value = res.http_peer
    up.value = res.up/1024
    down.value = res.down/1024
  })
}

const submitCallback = () => {
  if (newUrl.value !== "") {
    Add(newUrl.value).then(res => {
      console.log(res)
    })
  } else if (gameValue.value !== '' && httpValue.value !== '') {
    SetPeer(gameValue.value, httpValue.value).then(res => {
      console.log(res)
      getStatus()
    })
  } else if (newUrl.value !== '' && gameValue.value !== '' && httpValue.value !== '') {
    Add(newUrl.value).then(res => {
      console.log(res)
    })
    SetPeer(gameValue.value, httpValue.value).then(res => {
      console.log(res)
      getStatus()
    })
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
  margin-left: -50px;
}


.n-progress-content {
  width: 200px;
  height: 200px;
}


.n-progress-content svg {
  width: 200px;
  height: 200px;
}
</style>