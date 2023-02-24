<template>
    <div class="reasoning-service wrap">
        <div class="table-head">
            <a-input-search @search="onSearch" allowClear enterButton placeholder="请输入数据集名称"/>
            <a-button @click="$router.push({name: 'newInference'})" icon="plus" type="primary">新增</a-button>
        </div>
        <a-table :columns="columns" :dataSource="data" rowKey="ID">
            <template slot="framework" slot-scope="record">
          <span>
            {{record.Framework}}-{{record.FrameworkVersion}}
          </span>
            </template>
            <template slot="status" slot-scope="status">
                <span>
                    <a-tag :color="inferenceStatusTagColor(status)">{{inferenceStatusMap(status)}}</a-tag>
                    <a-icon v-if="status === 0" type="sync" :spin="true"/>
                </span>
            </template>
            <template slot="action" slot-scope="record">
                <span style="white-space: nowrap;">
<!--                    <a-button>停止</a-button>-->
<!--                    <a-divider type="vertical"/>-->
                    <a-button @click="_delete(record.ID)" type="danger">删除</a-button>
                </span>
            </template>
        </a-table>
    </div>
</template>

<script lang="ts">
    import {Component, Vue} from 'vue-property-decorator'
    import {deleteInference, getInference, getInferenceList, IInferenceListItem} from '@/api/inference'
    import {formatTime} from '@/utils/time-format'

    @Component
    export default class InferenceList extends Vue {
        private rowSelection = []
        private queryStatus: number = 0;

        private columns = [{
            title: '名称',
            dataIndex: 'Name'
        }, {
            title: '描述',
            dataIndex: 'Description'
        }, {
            title: '开发环境',
            scopedSlots: {customRender: 'framework'}
        }, {
            title: '创建时间',
            dataIndex: 'CreatedAt'
        }, {
            title: '状态',
            dataIndex: 'Status',
            scopedSlots: {customRender: 'status'}
        }, {
            title: '服务名',
            dataIndex: 'ServiceName'
        }, {
            title: '操作',
            scopedSlots: {customRender: 'action'}
        }]

        private data: IInferenceListItem[] = []

        private onSearch() {
        }

        private inferenceStatusMap(status: number): string {
            return ['创建中', '可访问', '失败'][status]
        }

        private inferenceStatusTagColor(status: number): string {
            return ['orange', 'green', 'red'][status]
        }

        private _delete(inferenceID: number): void {
            deleteInference(inferenceID).then(() => {
                this.$message.success('删除推理服务成功')
                this.data = this.data.filter((item) => {
                    return item.ID !== inferenceID
                })
            }).catch(() => {
                this.$message.error('删除推理服务失败')
            })
        }

        private created() {
            getInferenceList().then((res: any): void => {
                this.data = res.data.map((item: IInferenceListItem) => {
                    item.CreatedAt = formatTime(new Date(item.CreatedAt))
                    return item
                })
            })

            this.queryStatus = setInterval(() => {
                let unreadyList = this.data.filter((item) => {
                    return item.Status === 0
                })
                if (unreadyList.length === 0) {
                    clearInterval(this.queryStatus)
                    return
                }
                unreadyList.forEach((item) => {
                    getInference(item.ID).then((response) => {
                        item.Status = response.data.Status
                    })
                })
            }, 2000)
        }

        private beforeDestroy() {
            clearInterval(this.queryStatus)
        }
    }
</script>

<style lang="sass" scoped>
    .table-head
        display: flex
        margin: 8px 0
        justify-content: space-between

        .ant-input-search
            width: auto
</style>
