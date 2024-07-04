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
                  Game:{{ gamePeer === null ? '未选择' : gamePeer.name }}
                  <n-gradient-text v-if="gamePeer"
                                   :type="gamePeer.ping<60?'success':gamePeer.ping<100?'warning':'error'">
                    {{ gamePeer.ping }}
                  </n-gradient-text>
                </p>
                <p @click="getList()">
                  Http:{{ httpPeer === null ? '未选择' : httpPeer.name }}
                  <n-gradient-text v-if="httpPeer"
                                   :type="httpPeer.ping<60?'success':httpPeer.ping<100?'warning':'error'">
                    {{ httpPeer.ping }}
                  </n-gradient-text>
                </p>
              </n-space>
              <n-space vertical size="small" v-if="showUpDowInfo">
                <!--                <p>-->
                <!--                  上传:-->
                <!--                  <n-gradient-text v-if="up" type="success">-->
                <!--                    {{ up / 1024 > 1024 ? (up / 1024 / 1024).toFixed(2) + 'MB' : (up / 1024).toFixed(2) + 'KB' }}-->
                <!--                  </n-gradient-text>-->
                <!--                </p>-->
                <p>
                  流量统计:
                  <n-gradient-text v-if="down" type="success">
                    {{ down / 1024 > 1024 ? (down / 1024 / 1024).toFixed(2) + 'MB' : (down / 1024).toFixed(2) + 'KB' }}
                  </n-gradient-text>
                </p>
              </n-space>
            </n-space>
          </n-progress>
        </n-space>
        <n-space>
          <n-button :disabled="btnDisabled" @click="!state?start():stop()" style="margin-left: 110px">
            {{ btnText }}
          </n-button>
        </n-space>
        <n-gradient-text type="success" style="margin-left: 130px;margin-top: 35px">
          v1.2.8
        </n-gradient-text>
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
              :options="gameHttpOpt"
              placeholder="请选择Game"
              value-field="val"
              label-field="name"
          />
          <br>
          <n-select
              v-model:value="httpValue"
              vertical
              filterable
              :options="gameHttpOpt"
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
import {ref, defineComponent, Ref, reactive, onMounted, watch} from 'vue'
import {Add, List, SetPeer, Start, Status, Stop} from "../../wailsjs/go/main/App";
import {SelectOption, SelectGroupOption} from 'naive-ui'
import {onBeforeMount} from "@vue/runtime-core";
import {useMessage} from 'naive-ui'

const percentageRef = ref(0)
const state = ref(false)
const btnText = ref('开始加速')
const btnDisabled = ref(false)
const showModal = ref(false)
const gameHttpOpt = ref(Array<SelectOption | SelectGroupOption>())
const gameValue = ref()
const httpValue = ref()

const gamePeer: Ref<any> | null = ref(null)
const httpPeer: Ref<any> | null = ref(null)
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

const message = useMessage()

const start = () => {
  btnDisabled.value = true
  showGameHttpInfo.value = false
  showUpDowInfo.value = true
  btnText.value = '加速中.'
  Start().then(res => {
    if (res !== 'ok' && res !== 'running') {
      message.error(`加速失败:` + res)
      btnDisabled.value = false
      showUpDowInfo.value = false
      showGameHttpInfo.value = true
      return;
    }
    state.value = true
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
  })
}
const getList = () => {
  showModal.value = true
  gameHttpOpt.value = Array<SelectOption | SelectGroupOption>()
  List().then(res => {
    res.forEach((item) => {
      gameHttpOpt.value.push({
        name: item.name + '-' + item.ping + 'ms',
        val: item.name
      })
    })
  })
}


const getStatus = () => {
  Status().then(res => {
    if (res.game_peer !== null || res.http_peer !== null) {
      gamePeer.value = res.game_peer
      httpPeer.value = res.http_peer
      up.value = res.up
      down.value = res.down
      btnDisabled.value = false
      btnText.value = state.value ? '结束加速' : '开始加速'
      return;
    }
    btnText.value = '没有节点'
    btnDisabled.value = true
  })
}

const submitCallback = () => {
  if (newUrl.value !== undefined && gameValue.value !== undefined && httpValue.value !== undefined) {
    message.error('只能选择一种方式')
    newUrl.value = undefined;
    gameValue.value = undefined;
    httpValue.value = undefined;
    return
  }
  if (newUrl.value !== undefined) {
    Add(newUrl.value).then(res => {
      if (res === 'ok') {
        message.success('导入连接成功')
        newUrl.value = undefined;
      } else {
        message.error('导入连接失败')
      }
    });
  }
  if (gameValue.value !== undefined || httpValue.value !== undefined) {
    if (gameValue.value === undefined) {
      message.error('请选择Game节点')
      httpValue.value = undefined;
      return
    }
    if (httpValue.value === undefined) {
      message.error('请选择Http节点')
      gameValue.value = undefined;
      return
    }
    SetPeer(gameValue.value, httpValue.value).then(res => {
      if (res === 'ok') {
        message.success('设置节点成功')
        gameValue.value = undefined;
        httpValue.value = undefined;
      } else {
        message.error('设置节点失败')
      }
    });
  }
};


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
  margin-left: -100px;
}


.n-progress-content {
  width: 300px;
  height: 300px;
}


.n-progress-content svg {
  width: 300px;
  height: 300px;
}
</style>