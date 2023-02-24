<template>
  <div class="model-develop wrap">
    <div class="new-head">
      <h2>创建训练作业</h2>
      <a-button type="danger" @click="$router.push({name: 'jobList'})"> <a-icon type="left" />返回列表 </a-button>
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
        <a-form-item :label-col="{ span: 5 }" :wrapper-col="{ span: 12 }" label="训练代码">
            <a-select :filterOption="filterOption" @change="selectedCodeChange" optionFilterProp="children" placeholder="选择训练代码"
                      showSearch v-decorator="['CodeID', {
                rules: [{ required: true, message: '请选择训练代码!' }]
            }]">
                <template v-if="isFrameworkTensorFlow">
                    <a-select-option :key="c.ID" :value="c.ID" v-for="c in tensorFlowCodeList">{{c.Name}}</a-select-option>
                </template>
                <template v-else>
                    <a-select-option :key="c.ID" :value="c.ID" v-for="c in PYTorchCodeList">{{c.Name}}</a-select-option>
                </template>
            </a-select>
            <a-alert type="info">
                <template slot="message">
                    <div>代码存储路径：<strong>/workspace/code</strong></div>
                </template>
            </a-alert>
        </a-form-item>
        <template v-if="selectedCodeID !== 0">
            <a-form-item :label-col="{ span: 5 }" :wrapper-col="{ span: 12 }" label="入口文件"> <!--TODO: 更改为选择器-->
                <EntryPointSelector :codeID="selectedCodeID"
                                    v-decorator="['EntryPoint', { rules: [{ required: true, message: '请选择入口文件!' }] }]">
                </EntryPointSelector>
            </a-form-item>
            <a-form-item :label-col="{ span: 5 }" :wrapper-col="{ span: 12 }" label="运行参数">
                <ArgumentInput
                        v-decorator="['Args', {}]"
                />
            </a-form-item>
        </template>
        <a-form-item :label-col="{ span: 5 }" :wrapper-col="{ span: 12 }" label="训练数据集">
            <a-select :filterOption="filterOption" optionFilterProp="children" placeholder="选择训练数据集"
                      showSearch v-decorator="['DatasetID', {
                rules: [{ required: true, message: '请选择训练数据集!' }]
            }]">
                <a-select-option :key="d.ID" :value="d.ID" v-for="d in datasetList">{{d.Name}}</a-select-option>
            </a-select>
            <a-alert type="info">
                <template slot="message">
                    <div>数据集存储路径：<strong>/workspace/dataset</strong></div>
                </template>
            </a-alert>
        </a-form-item>
        <a-form-item :label-col="{ span: 5 }" :wrapper-col="{ span: 12 }" label="训练输出模型">
            <a-select :filterOption="filterOption" optionFilterProp="children" placeholder="选择训练输出模型"
                      showSearch v-decorator="['AIModelID', {
                rules: [{ required: true, message: '请选择训练输出模型!' }]
            }]">
                <template v-if="isFrameworkTensorFlow">
                    <a-select-option :key="m.ID" :value="m.ID" v-for="m in tensorFlowModelList">{{m.Name}}</a-select-option>
                </template>
                <template v-else>
                    <a-select-option :key="m.ID" :value="m.ID" v-for="m in PYTorchModelList">{{m.Name}}</a-select-option>
                </template>
            </a-select>
            <a-alert type="info">
                <template slot="message">
                    <div>模型存储路径：<strong>/workspace/model</strong></div>
                </template>
            </a-alert>
        </a-form-item>
        <a-form-item :label-col="{ span: 5 }" :wrapper-col="{ span: 12 }" label="节点数量">
            <a-slider :marks="NumMarks" :max="4" :min="1" class="slider" v-decorator="['Num', {
                rules: [{ required: true, message: '请选择节点数量!' }],
                initialValue: 1
            }]"/>
        </a-form-item>
        <a-form-item :label-col="{ span: 5 }" :wrapper-col="{ span: 12 }" label="处理器">
            <a-slider :marks="CPUMarks" :max="4" :min="1" class="slider" v-decorator="['CPU', {
                rules: [{ required: true, message: '请选择每个节点处理器数!' }],
                initialValue: 1
            }]"/>
        </a-form-item>
        <a-form-item :label-col="{ span: 5 }" :wrapper-col="{ span: 12 }" label="内存">
            <a-slider :marks="MemoryMarks" :max="4" :min="1" class="slider" v-decorator="['Memory', {
                rules: [{ required: true, message: '请选择每个节点内存!' }],
                initialValue: 1
            }]"/>
        </a-form-item>
        <a-form-item :label-col="{ span: 5 }" :wrapper-col="{ span: 12 }" label="GPU">
            <a-slider :marks="GPUMarks" :max="2" :min="1" class="slider" v-decorator="['GPU', {
                rules: [{ required: true, message: '请选择每个节点GPU数!' }],
                initialValue: 1
            }]"/>
        </a-form-item>
        <a-form-item :wrapper-col="{ span: 12, offset: 5 }">
            <a-button html-type="submit" type="primary" :loading="isSubmitting">
                确认创建训练任务
            </a-button>
        </a-form-item>
    </a-form>
  </div>
</template>

<script lang="ts">
    import {Component, Vue} from 'vue-property-decorator'
    import ArgumentInput from '@/components/ArgumentInput.vue'
    import FrameworkSelector from '@/components/framework/FrameworkSelector.vue'
    import FrameworkVersionSelector from '@/components/framework/FrameworkVersionSelector.vue'
    import {getCodeList, ICodeListItem} from '@/api/code'
    import {getDatasetList, IDatasetListItem} from '@/api/dataset'
    import {getModelList, IModelListItem} from '@/api/model'
    import EntryPointSelector from '@/components/EntryPointSelector.vue'
    import {createNewJob} from '@/api/jobs'

    @Component({
        components: {EntryPointSelector, FrameworkVersionSelector, FrameworkSelector, ArgumentInput}
    })
    export default class NewJob extends Vue {
        private form!: any;
        private isFrameworkTensorFlow: boolean = true;
        private selectedCodeID: number = 0;
        private NumMarks = {1: '1节点', 2: '2节点', 3: '3节点', 4: '4节点'};
        private CPUMarks = {1: '1核', 2: '2核', 3: '3核', 4: '4核'};
        private MemoryMarks = {1: '1Gi', 2: '2Gi', 3: '3Gi', 4: '4Gi'};
        private GPUMarks = {1: '1块', 2: '2块'};
        private isSubmitting: boolean = false;

        // TODO: 获取code、dataset、model列表
        private codeList: ICodeListItem[] = [];
        private datasetList: IDatasetListItem[] = [];
        private modelList: IModelListItem[] = [];

        private get tensorFlowCodeList(): ICodeListItem[] {
            return this.codeList.filter((item) => {
                return item.Framework === 'tensorflow'
            })
        }

        private get PYTorchCodeList(): ICodeListItem[] {
            return this.codeList.filter((item) => {
                return item.Framework === 'pytorch'
            })
        }

        private get tensorFlowModelList(): IModelListItem[] {
            return this.modelList.filter((item) => {
                return item.Framework === 'tensorflow'
            })
        }

        private get PYTorchModelList(): IModelListItem[] {
            return this.modelList.filter((item) => {
                return item.Framework === 'pytorch'
            })
        }

        private handleSubmit(e: any) {
            e.preventDefault()
            this.form.validateFields((err: any, values: any) => {
                if (!err) {
                    this.isSubmitting = true
                    createNewJob(values).then((response) => {
                        this.$message.success('创建训练任务成功')
                        this.$router.replace({
                            name: 'jobList'
                        })
                    }).catch(() => {
                        this.$message.error('创建训练任务失败，请重试')
                    }).finally(() => {
                        this.isSubmitting = false
                    })
                }
            })
        }

        private selectedCodeChange(value: number): void {
            this.selectedCodeID = value
        }

        private frameworkChange(value: string): void {
            this.isFrameworkTensorFlow = value === 'tensorflow'
            this.selectedCodeID = 0
            this.form.resetFields(['FrameworkVersion'])
            this.form.resetFields(['CodeID'])
            this.form.resetFields(['AIModelID'])
            this.form.resetFields(['EntryPoint'])
        }

        private filterOption(input: string, option: any): any {
            return option.componentOptions.children[0].text.toLowerCase().includes(input.toLowerCase())
        }

        private created(): void {
            this.form = this.$form.createForm(this)
            getCodeList().then((response) => {
                this.codeList = response.data
            })
            getDatasetList().then((response) => {
                this.datasetList = response.data
            })
            getModelList().then((response) => {
                this.modelList = response.data
            })
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
>>>.slider
    .ant-slider-rail
        background-color: #c3c3c3
    .ant-slider-mark-text
        white-space: nowrap
</style>
