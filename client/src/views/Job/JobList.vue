<template>
    <div class="model-train wrap">
        <div class="table-head">
            <a-input-search @search="onSearch" allowClear enterButton placeholder="请输入数据集名称"/>
            <a-button @click="$router.push({name: 'newJob'})" icon="plus" type="primary">新增</a-button>
        </div>
        <a-table :columns="columns" :dataSource="data" rowKey="ID">
            <template slot="framework" slot-scope="record">
                <span>
                    {{record.Framework}}-{{record.FrameworkVersion}}
                </span>
            </template>
            <template slot="status" slot-scope="status">
                <span>
                    <a-tag :color="jobStatusTagColor(status)">{{jobStatusMap(status)}}</a-tag>
                    <a-icon v-if="status === 0 || status === 1" type="sync" :spin="true"/>
                </span>
            </template>
            <template slot="action" slot-scope="record">
                <span style="white-space: nowrap;">
                    <a-button>停止</a-button>
                    <a-divider type="vertical"/>
                    <a-button @click="_delete(record.ID)" type="danger">删除</a-button>
                </span>
            </template>
        </a-table>
    </div>
</template>

<script lang="ts">
    import {Component, Vue} from 'vue-property-decorator'
    import {getModelList} from '@/api/model'
    import {deleteJob, getJobsList, IJobListItem, getJob} from '@/api/jobs'
    import {formatTime} from '@/utils/time-format'

    @Component
    export default class JobList extends Vue {
        private rowSelection = []
        private queryStatus: number = 0

        private columns = [{
            title: '名称',
            dataIndex: 'Name'
        }, {
            title: '描述',
            dataIndex: 'Description'
        }, {
            title: '运行环境',
            scopedSlots: {customRender: 'framework'}
        }, {
            title: '创建时间',
            dataIndex: 'CreatedAt'
        }, {
            title: '状态',
            dataIndex: 'Status',
            scopedSlots: {customRender: 'status'}
        }, {
            title: '操作',
            scopedSlots: {customRender: 'action'}
        }]

        private data: IJobListItem[] = []

        private onSearch() {
        }

        private jobStatusMap(status: number): string {
            return ['创建中', '训练中', '已完成', '已停止', '失败'][status]
        }

        private jobStatusTagColor(status: number): string {
            return ['orange', 'cyan', 'green', 'blue', 'red'][status]
        }

        private _delete(jobID: number): void {
            deleteJob(jobID).then(() => {
                this.$message.success('删除训练任务成功')
                this.data = this.data.filter((item) => {
                    return item.ID !== jobID
                })
            }).catch(() => {
                this.$message.error('删除训练任务失败')
            })
        }

        private created() {
            getJobsList().then((res: any): void => {
                this.data = res.data.map((item: IJobListItem) => {
                    item.CreatedAt = formatTime(new Date(item.CreatedAt))
                    return item
                })
            })

            this.queryStatus = setInterval(() => {
                let unreadyList = this.data.filter((item) => {
                    return item.Status === 0 || item.Status === 1
                })
                if (unreadyList.length === 0) {
                    clearInterval(this.queryStatus)
                    return
                }
                unreadyList.forEach((item) => {
                    getJob(item.ID).then((response) => {
                        item.Status = response.data.Status
                    })
                })
            }, 5000)
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
