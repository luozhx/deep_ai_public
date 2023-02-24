<template>
    <div class="console">
        <div class="process">
            <a-row>
                <a-col :span="4" class="process_img" offset="4"></a-col>
                <a-col :span="4" class="process_img"></a-col>
                <a-col :span="4" class="process_img"></a-col>
                <a-col :span="4" class="process_img"></a-col>
            </a-row>
            <a-row>
                <a-col :span="4" class="process_name" offset="4">
                    <span class="icon_num">1</span>代码开发
                </a-col>
                <a-col :span="4" class="process_name">
                    <span class="icon_num">2</span>数据管理
                </a-col>
                <a-col :span="4" class="process_name">
                    <span class="icon_num">3</span>模型管理
                </a-col>
                <a-col :span="4" class="process_name">
                    <span class="icon_num">4</span>训练结果
                </a-col>
            </a-row>
        </div>
        <div class="charts">
            <a-row type="flex" justify="center">
                <a-col :span="10" class="charts-block">
                    <div ref="code"></div>
                </a-col>
                <a-col :span="10" class="charts-block">
                    <div ref="model"></div>
                </a-col>
            </a-row>
            <a-row type="flex" justify="center">
                <a-col :span="10" class="charts-block">
                    <div ref="job"></div>
                </a-col>
                <a-col :span="10" class="charts-block">
                    <div ref="inference"></div>
                </a-col>
            </a-row>
        </div>
    </div>
</template>

<script lang="ts">
    import {Component, Vue} from 'vue-property-decorator'
    import PieCharts from '@/components/PieCharts.vue'
    import {Ring} from '@antv/g2plot'
    import {getSystemStatus} from '@/api/console'

    @Component({
        components: {
            PieCharts
        }
    })
    export default class Console extends Vue {
        mounted() {
            getSystemStatus().then((res: any): void => {
                for (let k in res.data) {
                    if (k === 'dataset') continue

                    new Ring(this.$refs[k] as HTMLDivElement, {
                        forceFit: true,
                        title: {
                            visible: true,
                            text: ({code: '代码', model: '模型', job: '训练任务', inference: '推理服务'} as any)[k]
                        },
                        radius: 0.8,
                        padding: 'auto',
                        data: res.data[k].map((item: any) => {
                            item.status = this.codeStatusMap(k, item.status)
                            return item
                        }),
                        angleField: 'count',
                        colorField: 'status'
                    }).render()
                }
            })
        }

        private codeStatusMap(type: string, status: number): string {
            switch (type) {
                case 'code':
                    return ['不可用', '可用', '被占用'][status]
                case 'model':
                    return ['创建中', '不可用', '可用', '被占用'][status]
                case 'job':
                    return ['创建中', '训练中', '已完成', '已停止', '失败'][status]
                case 'inference':
                    return ['创建中', '可访问', '失败'][status]
                default:
                    return ''
            }
        }
    }
</script>

<style lang="scss" scoped>
    .process_img {
        height: 100px;
        background: url(../../public/images/icon_debug.png) center no-repeat;
        background-size: 78px 78px;
    }

    .process_name {
        display: flex;
        justify-content: center;
    }

    .icon_num {
        display: flex;
        justify-content: center;
        align-items: center;
        width: 20px;
        height: 20px;
        color: white;
        background-color: rgb(173, 195, 192);
    }

    .charts {
        margin-top: 50px;
    }

    .charts-block {
        background-color: white;
        border: 1px dashed #ccc;
    }
</style>
