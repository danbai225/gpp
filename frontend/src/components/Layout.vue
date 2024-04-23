<template>
  <n-space vertical>
<!--    <n-space>-->
<!--      <n-switch v-model:value="inverted"/>-->
<!--      inverted-->
<!--    </n-space>-->

    <n-layout>
      <n-layout has-sider>
        <n-layout-sider
            class="height"
            bordered
            show-trigger
            collapse-mode="width"
            :collapsed-width="64"
            :width="240"
            :native-scrollbar="false"
            :inverted="inverted"
            style="max-height: 320px"
        >
          <n-menu
              :inverted="inverted"
              :collapsed-width="64"
              :collapsed-icon-size="22"
              :options="menuOptions"
          />
        </n-layout-sider>
        <n-layout style="max-height: 320px"/>
        <n-layout/>
      </n-layout>
<!--      <n-layout-footer :inverted="inverted" bordered>-->
<!--        Footer Footer Footer-->
<!--      </n-layout-footer>-->
    </n-layout>
  </n-space>
</template>

<script lang="ts">
import {h, defineComponent, ref, Component} from 'vue'
import {NIcon} from 'naive-ui'
import {
  BookOutline as BookIcon, GameController, List,
  PersonOutline as PersonIcon,
  WineOutline as WineIcon
} from '@vicons/ionicons5'

function renderIcon(icon: Component) {
  return () => h(NIcon, null, {default: () => h(icon)})
}

const menuOptions = [
  {
    label: '游戏列表',
    key: 'game-list',
    icon: renderIcon(List)
  },
  {
    label: '全部游戏',
    key: 'all-games',
    icon: renderIcon(GameController),
  },
  {
    label: '全部游戏',
    key: 'all-games',
    icon: renderIcon(GameController),
  },
  {
    label: '全部游戏',
    key: 'all-games',
    icon: renderIcon(GameController),
  },
  {
    label: '全部游戏',
    key: 'all-games',
    icon: renderIcon(GameController),
  },
]

export default defineComponent({
  setup() {
    return {
      inverted: ref(false),
      menuOptions
    }
  }
})
</script>

<style>
</style>