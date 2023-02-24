<template>
    <div>
        <div :key="index" class="argument-item" v-for="(a, index) in trainingArguments">
            <div class="label">参数{{index + 1}}</div>
            <a-input @change="$emit('change', validTrainingArguments)" class="argument-item-input"
                     v-model="a.key"></a-input>
            <div>=</div>
            <a-input @change="$emit('change', validTrainingArguments)" class="argument-item-input"
                     v-model="a.value"></a-input>
            <a-button @click="delArgument(index)" icon="delete" type="danger"></a-button>
        </div>
        <div class="arguments-string-box" v-if="trainingArguments.length !== 0">
            <div class="label">结果参数字符串：</div>
            <div class="arguments-string">{{validTrainingArguments.join(' ')}}</div>
        </div>
        <a-button @click="addArgument" icon="plus" type="primary"></a-button>
    </div>
</template>

<script lang="ts">

    import {Component, Prop, Vue} from 'vue-property-decorator'

    interface ArgumentItem {
        key: string;
        value: string;
    }

    @Component
    export default class ArgumentInput extends Vue {
        @Prop({type: Array}) private readonly value!: string[];

        private trainingArguments: ArgumentItem[] = [];

        private get validTrainingArguments(): string[] {
            return this.trainingArguments.filter((item) => {
                return item.key.trim() !== '' && item.value.trim() !== ''
            }).map((item) => {
                return `--${item.key}=${item.value}`
            })
        }

        private addArgument(): void {
            this.trainingArguments.push({
                key: '',
                value: ''
            } as ArgumentItem)
        }

        private delArgument(idx: number) {
            this.trainingArguments = this.trainingArguments.filter((item, index) => {
                return index !== idx
            })
            this.$nextTick(() => {
                this.$emit('change', this.trainingArguments)
            })
        }

        private created(): void {
            if (this.value === undefined) {
                this.$emit('change', [])
                return
            }
            const newArguments: ArgumentItem[] = []
            this.value.forEach((item) => {
                const [key, value] = item.split('=')
                newArguments.push({key, value})
            })
            this.trainingArguments = newArguments
        }
    }
</script>

<style scoped>
    .argument-item {
        display: flex;
        align-items: center;
        flex-direction: row;
        justify-content: space-between;
        margin-bottom: 10px;
    }

    .label {
        white-space: nowrap;
        font-weight: bold;
    }

    .arguments-string-box {
        display: flex;
        align-items: center;
        flex-direction: row;
        margin-bottom: 10px;
    }

    .argument-item-input {
        margin: 0 10px;
    }

    .arguments-string {
        width: 100%;
        line-height: 20px;
        font-weight: bold;
        font-size: 16px;
    }
</style>
