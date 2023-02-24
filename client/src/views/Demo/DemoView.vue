<template>
  <a-row :gutter="16">
    <a-col :span="4" v-for="(item, index) in demoList" :key="index" style="margin-bottom: 12px">
       <a-card
        hoverable
        @click="$router.push({name:item.redirect})"
       >
        <img
          :src="item.img"
          slot="cover"
        />
        <a-card-meta :title="item.title" />
      </a-card>
    </a-col>
  </a-row>
</template>

<script lang="ts">
import {Component, Vue} from 'vue-property-decorator'

@Component
export default class extends Vue {
  private demoList = [{
    img: '/images/demo/minst.jpg',
    title: 'minst手写数据集demo',
    redirect: 'demoMnist'
  }, {
    img: '/images/demo/edsr.png',
    title: '超分辨率',
    redirect: 'demoEDSR'
  }]
}
</script>
