<template>
    <div class="model-develop wrap">
        <div class="table-head">
            <a-input-search @search="onSearch" allowClear enterButton placeholder="请输入数据集名称"/>
            <a-button @click="$router.push({name: 'newCode'})" icon="plus" type="primary">新增</a-button>
        </div>
        <a-table :columns="columns" :dataSource="data" rowKey="ID">
            <template slot="framework" slot-scope="record">
                <span>{{record.Framework}}-{{record.FrameworkVersion}}</span>
            </template>
            <template slot="PVCStatus" slot-scope="PVCStatus">
                <span>
                    <a-tag :color="PVCStatusTagColor(PVCStatus)">{{PVCStatusMap(PVCStatus)}}</a-tag>
                    <a-icon v-if="PVCStatus === 0" type="sync" :spin="true"/>
                </span>
            </template>
            <template slot="notebookStatus" slot-scope="notebookStatus">
                <span>
                    <a-tag :color="notebookStatusTagColor(notebookStatus)">{{notebookStatusMap(notebookStatus)}}</a-tag>
                    <a-icon v-if="notebookStatus === 0" type="sync" :spin="true"/>
                </span>
            </template>
            <template slot="status" slot-scope="status">
                <span>
                    <a-tag :color="codeStatusTagColor(status)">{{codeStatusMap(status)}}</a-tag>
                </span>
            </template>
            <template slot="action" slot-scope="record">
                <span style="white-space: nowrap;">
                    <a-button type="primary" @click="openNotebook(record.ID, record.ServiceName)"
                              :disabled="record.NotebookStatus !== 1">打开Notebook</a-button>
                    <a-divider type="vertical"/>
                    <a-button @click="download(record)" :disabled="record.PVCStatus !== 1">下载</a-button>
                    <a-divider type="vertical"/>
                    <a-button @click="_delete(record.ID)" type="danger">删除</a-button>
                </span>
            </template>
        </a-table>
    </div>
</template>

<script lang="ts">
    import {Component, Vue} from 'vue-property-decorator'
    import {deleteCode, getCode, getCodeDownloadUrl, getCodeList, getCodeNotebookUrl, ICodeListItem} from '@/api/code'
    import {formatTime} from '@/utils/time-format'

    @Component
    export default class CodeList extends Vue {
        private rowSelection = []
        private queryStatus: number = 0

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
            title: '存储资源状态',
            dataIndex: 'PVCStatus',
            scopedSlots: {customRender: 'PVCStatus'}
        }, {
            title: 'Notebook状态',
            dataIndex: 'NotebookStatus',
            scopedSlots: {customRender: 'notebookStatus'}
        }, {
            title: '状态',
            dataIndex: 'Status',
            scopedSlots: {customRender: 'status'}
        }, {
            title: '操作',
            scopedSlots: {customRender: 'action'}
        }]

        private data: ICodeListItem[] = [];

        private onSearch() {
        }

        private PVCStatusMap(status: number): string {
            return ['创建中', '正常', '不可用'][status]
        }

        private PVCStatusTagColor(status: number): string {
            return ['orange', 'green', 'red'][status]
        }

        private notebookStatusMap(status: number): string {
            return ['创建中', '可访问', '失败'][status]
        }

        private notebookStatusTagColor(status: number): string {
            return ['orange', 'green', 'red'][status]
        }

        private codeStatusMap(status: number): string {
            return ['不可用', '可用', '被占用'][status]
        }

        private codeStatusTagColor(status: number): string {
            return ['red', 'green', 'cyan'][status]
        }

        private _delete(codeID: number): void {
            deleteCode(codeID).then(() => {
                this.$message.success('删除代码成功')
                this.data = this.data.filter((item) => {
                    return item.ID !== codeID
                })
            }).catch(() => {
                this.$message.error('删除代码失败')
            })
        }

        private openNotebook(codeID: number, serviceName: string): void {
            const link = document.createElement('a')
            link.href = getCodeNotebookUrl(codeID, serviceName)
            link.target = '_blank'
            document.body.appendChild(link)
            link.click()
            document.body.removeChild(link)
        }

        private download(code: ICodeListItem): void {
            const link = document.createElement('a')
            link.href = getCodeDownloadUrl(code.ID)
            link.download = `${code.Name}.zip`
            document.body.appendChild(link)
            link.click()
            document.body.removeChild(link)
        }

        private created() {
            getCodeList().then((res: any): void => {
                this.data = res.data.map((item: ICodeListItem) => {
                    item.CreatedAt = formatTime(new Date(item.CreatedAt))
                    return item
                })
            })

            this.queryStatus = setInterval(() => {
                let unreadyList = this.data.filter((item) => {
                    return item.PVCStatus === 0 || item.NotebookStatus === 0
                })
                if (unreadyList.length === 0) {
                    clearInterval(this.queryStatus)
                    return
                }
                unreadyList.forEach((item) => {
                    getCode(item.ID).then((response) => {
                        [item.PVCStatus, item.NotebookStatus, item.Status] =
                            [response.data.PVCStatus, response.data.NotebookStatus, response.data.Status]
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
