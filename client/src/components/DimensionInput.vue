<template>
    <div>
        <div style="display: flex;flex-direction: row;height: 40px;align-items: center;">
            <a-button type="primary" icon="minus" class="minus-btn" @click="minus"></a-button>
            <a-input class="number-input" :value="dimArray.length" :readOnly="true"/>
            <a-button type="primary" icon="plus" class="plus-btn" @click="plus"></a-button>
        </div>
        <div style="display: flex;flex-direction: row;flex-wrap: wrap;margin-top: 15px" v-if="dimArray.length > 0">
            <template v-for="(d, index) in dimArray">
                <div :key="index" >
                    <strong>第{{index + 1}}维</strong><a-input class="input-item" @input="onChange(index, $event)" :value="d"></a-input>
                </div>
            </template>
        </div>
    </div>
</template>

<script lang="ts">
    import {Component, Prop, Vue, Watch} from 'vue-property-decorator'

    @Component
    export default class DimensionInput extends Vue {
        @Prop({type: String}) private readonly value!: string;
        private dimArray: number[] = [];

        @Watch('value')
        private updateDimArray(): void {
            if (!this.value) {
                this.$emit('change', '')
                return
            }
            this.dimArray = this.value.split(' ').map((item) => {
                return parseInt(item, 10)
            })
        }

        private onChange(index: number, e: Event): void {
            this.$set(this.dimArray, index, !(e.target as any).value ? 0 : parseInt((e.target as any).value, 10))
            this.$emit('change', this.dimArray.join(' '))
        }

        private created(): void {
            this.updateDimArray()
        }

        private minus(): void {
            this.dimArray.pop()
            this.$emit('change', this.dimArray.join(' '))
        }

        private plus(): void {
            this.dimArray.push(0)
            this.$emit('change', this.dimArray.join(' '))
        }
    }
</script>

<style scoped>
    .minus-btn {
        border-radius: 4px 0 0 4px;
    }

    .number-input {
        border-radius: 0;
        text-align: center;
        width: 60px;
    }

    .plus-btn {
        border-radius: 0 4px 4px 0;
    }

    .input-item:first-child{
        margin-left: 0;
    }

    .input-item {
        text-align: center;
        margin: 0 10px;
        width: 60px;
    }
</style>
