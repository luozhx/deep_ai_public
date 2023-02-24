<template>
  <div class="model-develop wrap">
    <div class="new-head">
      <h2>创建开发环境</h2>
      <a-button type="danger" @click="$router.replace({name: 'codeList'})"> <a-icon type="left" />返回列表 </a-button>
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
        <a-form-item :label-col="{ span: 5 }" :wrapper-col="{ span: 12 }" label="环境">
            <FrameworkSelector @change="frameworkChange" v-decorator="['Framework', {
                    rules: [{ required: true, message: '请选择环境!' }],
                    initialValue: 'tensorflow'
                }]">
            </FrameworkSelector>
        </a-form-item>
        <a-form-item :label-col="{ span: 5 }" :wrapper-col="{ span: 12 }" label="环境版本">
            <FrameworkVersionSelector :type="isFrameworkTensorFlow?'tensorflow':'pytorch'" v-decorator="['FrameworkVersion', {
                    rules: [{ required: true, message: '请选择环境版本!' }],
                    initialValue: isFrameworkTensorFlow?'1.15.0':'1.3'
                }]">
            </FrameworkVersionSelector>
        </a-form-item>
        <a-form-item :wrapper-col="{ span: 12, offset: 5 }">
            <a-button html-type="submit" type="primary" :loading="isSubmitting">
                确认创建代码环境
            </a-button>
        </a-form-item>
    </a-form>
  </div>
</template>

<script lang="ts">
    import {Component, Vue} from 'vue-property-decorator'
    import FrameworkSelector from '@/components/framework/FrameworkSelector.vue'
    import FrameworkVersionSelector from '@/components/framework/FrameworkVersionSelector.vue'
    import {createNewCode} from '@/api/code'

    @Component({
        components: {FrameworkVersionSelector, FrameworkSelector}
    })
    export default class NewCode extends Vue {
        private form!: any;
        private isFrameworkTensorFlow: boolean = true;
        private isSubmitting: boolean = false;

        private handleSubmit(e: any) {
            e.preventDefault()
            this.form.validateFields((err: any, values: any) => {
                if (!err) {
                    this.isSubmitting = true
                    createNewCode(values).then((response) => {
                        this.$message.success('创建代码环境成功')
                        this.$router.replace({
                            name: 'codeList'
                        })
                    }).catch(() => {
                        this.$message.error('创建代码环境失败，请重试')
                    }).finally(() => {
                        this.isSubmitting = false
                    })
                }
            })
        }

        private frameworkChange(value: string): void {
            this.isFrameworkTensorFlow = value === 'tensorflow'
            this.form.resetFields(['FrameworkVersion'])
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
