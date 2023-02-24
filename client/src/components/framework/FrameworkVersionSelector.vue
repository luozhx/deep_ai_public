<template>
    <a-radio-group :value="value" @change="onChange" buttonStyle="solid">
        <template v-if="type === 'tensorflow' ">
            <a-radio-button value="1.15.0">1.15.0</a-radio-button>
            <a-radio-button value="1.14.0">1.14.0</a-radio-button>
        </template>
        <template v-else>
            <a-radio-button value="1.3">1.3</a-radio-button>
            <a-radio-button value="1.2">1.2</a-radio-button>
        </template>
    </a-radio-group>
</template>

<script lang="ts">
    import {Component, Prop, Vue} from 'vue-property-decorator'

    @Component
    export default class FrameworkVersionSelector extends Vue {
        @Prop({type: String}) private readonly value!: string;
        @Prop({type: String}) private readonly type!: string;

        private onChange(e: Event): void {
            this.$emit('change', (e.target as any).value)
        }
    }
</script>

<style scoped>

</style>
