<template>
    <div>
        <a-button @click="onClick">{{!value || value === '' ? '选择入口文件':`已选择: ${value}`}}</a-button>
        <a-modal @ok="confirmSelection" cancelText="取消" okText="确认" title="选择入口文件" v-model="visible">
            <a-directory-tree :loadData="onLoadData" :selectedKeys="[selection]" :treeData="data"
                              @select="fileSelected">
            </a-directory-tree>
        </a-modal>
    </div>
</template>

<script lang="ts">
    import {Component, Prop, Vue, Watch} from 'vue-property-decorator'
    import {getCodeFileList, ICodeFileListItem} from '@/api/code'

    @Component
    export default class EntryPointSelector extends Vue {
        @Prop({type: String}) private readonly value!: string;
        @Prop({type: Number, required: true}) private readonly codeID!: number;

        @Watch('value')
        private updateSelection(): void {
            this.selection = this.value
        }

        @Watch('codeID')
        private updateCodeDir(): void {
            getCodeFileList(this.codeID, '').then((response) => {
                this.data = response.data.map((item: ICodeFileListItem) => {
                    return {
                        title: item.name,
                        key: item.name,
                        isLeaf: !item.isDir,
                        selectable: !item.isDir
                    }
                })
            })
        }

        private data: any = [];
        private selection: string = '';
        private visible: boolean = false;

        private created(): void {
            if (this.value) {
                this.selection = this.value
            }
            this.updateCodeDir()
        }

        private onClick(): void {
            this.visible = true
            this.updateCodeDir()
        }

        private confirmSelection(): void {
            if (this.selection === '') {
                return
            }
            this.visible = false
            this.$emit('change', this.selection)
        }

        private fileSelected(selectedKeys: string[]): void {
            if (selectedKeys.length !== 0) {
                this.selection = selectedKeys[0]
            }
        }

        private onLoadData(treeNode: any): Promise<void> {
            return new Promise(resolve => {
                if (treeNode.dataRef.children) {
                    resolve()
                    return
                }
                getCodeFileList(this.codeID, `${treeNode.dataRef.key}`).then((response) => {
                    treeNode.dataRef.children = response.data.map((item: ICodeFileListItem) => {
                        return {
                            title: item.name,
                            key: `${treeNode.dataRef.key}/${item.name}`,
                            isLeaf: !item.isDir,
                            checkable: !item.isDir
                        }
                    })
                    this.data = [...this.data]
                    resolve()
                })
            })
        }
    }
</script>

<style scoped>

</style>
