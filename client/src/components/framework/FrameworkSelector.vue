<template>
    <a-radio-group :value="value" @change="onChange" buttonStyle="solid">
        <a-radio-button value="tensorflow">TensorFlow</a-radio-button>
        <a-radio-button value="pytorch">PYTorch</a-radio-button>
    </a-radio-group>
</template>

<script lang="ts">

    import {Component, Prop, Vue} from 'vue-property-decorator'

    @Component
    export default class FrameworkSelector extends Vue {
        @Prop({type: String}) private readonly value!: string;

        private onChange(e: Event): void {
            this.$emit('change', (e.target as any).value)
        }
    }
</script>

<style scoped>

</style>
