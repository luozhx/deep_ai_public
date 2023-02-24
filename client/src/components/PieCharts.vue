<template>
  <div class="console">
    <!-- <div :id="this.chartsId" style="width: 550px;height: 255px;margin-top: 5px;"></div> -->
    <v-chart theme="normal" :options="options"/>
  </div>
</template>

<script lang="ts">
import {Component, Prop, Vue} from 'vue-property-decorator'

@Component
export default class CircleCharts extends Vue {
  @Prop({
    default: '测试数据'
  })
  private chartsTitle!:String

  @Prop({
    default: () => [
      {value: 335, name: 'tensorflow-cpu'},
      {value: 310, name: 'tf-gpu'},
      {value: 234, name: 'picar'},
      {value: 135, name: 'vnc-test'},
      {value: 1548, name: 'jupyter-proxy-test'}
    ]
  })
  private readonly data!: Array<{
    value: Number,
    name: String
  }>
  // private test: string = 'id'
  // public $echarts: any
  private options = {
    title: {
      text: this.chartsTitle,
      x: 'left'
    },
    toolip: {
      trigger: 'item',
      formatter: '{a} <br/>{b} : {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      right: 30,
      top: 20,
      bottom: 20,
      data: this.data.map(item => item.name)
      // data: [
      //   'tensorflow-cpu',
      //   'tf-gpu',
      //   'picar',
      //   'vnc-test',
      //   'jupyter-proxy-test'
      // ]
    },
    // color: ['#F74461', '#576069', '#A9A9A9', '#ADC3C0', '#ADD8E6'],
    // stillShowZeroSum: false,
    series: [
      {
        name: 'XXXX',
        type: 'pie',
        radius: '55%',
        center: ['50%', '50%'],
        data: this.data,
        // data: [
        //   {value: 335, name: 'tensorflow-cpu'},
        //   {value: 310, name: 'tf-gpu'},
        //   {value: 234, name: 'picar'},
        //   {value: 135, name: 'vnc-test'},
        //   {value: 1548, name: 'jupyter-proxy-test'}
        // ],
        itemStyle: {
          emphasis: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        }
      }
    ]
  }
}
</script>

<style lang="sass" scoped>
.console
  height: 280px
  .echarts
    width: 100%
    height: 100%
</style>
