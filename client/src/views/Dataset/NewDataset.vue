<template>
    <div class="model-develop wrap">
        <div class="new-head">
        <h2>创建数据集</h2>
        <a-button type="danger" @click="$router.push({name: 'datasetList'})"> <a-icon type="left" />返回列表 </a-button>
        </div>
        <a-form :form="form" @submit="handleSubmit">
            <a-form-item :label-col="{ span: 5 }" :wrapper-col="{ span: 12 }" label="名称">
                <a-input
                        v-decorator="['Name', { rules: [{ required: true, message: '请输入名称!' }] }]"
                />
            </a-form-item>
            <a-form-item :label-col="{ span: 5 }" :wrapper-col="{ span: 12 }" label="描述">
                <a-input
                        v-decorator="['Description', { rules: [{ required: true, message: '请输入描述!' }] }]"
                />
            </a-form-item>
            <a-form-item :label-col="{ span: 5 }" :wrapper-col="{ span: 12 }" label="数据集文件">
                <a-upload :beforeUpload="beforeUpload"
                        v-decorator='["file", {
                                valuePropName: "fileList",
                                getValueFromEvent: normFile,
                                rules: [
                                    {required: true, message: "请选择数据集文件"},
                                ],
                                initialValue: null,
                                }]'>
                    <a-button>
                        <a-icon type="upload"/>
                        选择数据集文件
                    </a-button>
                </a-upload>
            </a-form-item>
            <a-form-item :wrapper-col="{ span: 12, offset: 5 }">
                <a-button html-type="submit" type="primary" :loading="isSubmitting">
                    确认创建数据集
                </a-button>
            </a-form-item>
        </a-form>
    </div>
</template>

<script lang="ts">
    import {Component, Vue} from 'vue-property-decorator'
    import {createNewDataset} from '@/api/dataset'

    @Component
    export default class NewDataset extends Vue {
        private form!: any;
        private isSubmitting: boolean = false;

        private handleSubmit(e: any) {
            e.preventDefault()
            this.form.validateFields((err: any, values: any) => {
                if (!err) {
                    this.isSubmitting = true

                    const formData = new FormData()
                    for (const key of Object.keys(values)) {
                        formData.append(key, key === 'file' ? values[key][0].originFileObj : values[key])
                    }

                    createNewDataset(formData).then((response) => {
                        this.$message.success('创建数据集成功')
                        this.$router.replace({
                            name: 'datasetList'
                        })
                    }).catch(() => {
                        this.$message.error('创建数据集失败，请重试')
                    }).finally(() => {
                        this.isSubmitting = false
                    })
                }
            })
        }

        private normFile({fileList}: any): any {
            return fileList.slice(-1)
        }

        private beforeUpload(): boolean {
            return false
        }

        private created(): void {
            this.form = this.$form.createForm(this)
        }
    }
</script>

<style lang="sass" scoped>
.new-head
  display: flex
  margin: 8px 0
  justify-content: space-between
  .ant-input-search
    width: auto

</style>
