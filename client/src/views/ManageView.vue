<template>
    <a-layout id="mainLayout">
        <a-layout>
            <a-layout-sider
                    style="background: #fff"
                    width="250"
            >
                <div class="slider-head">
                    <div class="logo"></div>
                    <h2>深度学习解决方案</h2>
                </div>
                <a-menu
                        :defaultOpenKeys="['sub1']"
                        :defaultSelectedKeys="['1']"
                        :style="{ borderRight: 0 }"
                        mode="inline"
                >
                    <a-menu-item :key="index" v-for="(item, index) in navList">
                        <a-icon :type="item.icon"/>
                        <router-link :to="`/${item.path}`" style="display: inline-block">{{item.name}}</router-link>
                    </a-menu-item>
                </a-menu>
            </a-layout-sider>
            <a-layout style="padding: 0 24px 24px">
                <a-breadcrumb style="margin: 16px 0">
                    <a-breadcrumb-item :key="index" v-for="(item, index) in breadcrumbList">{{item}}</a-breadcrumb-item>
                </a-breadcrumb>
                <a-layout-content>
                    <router-view/>
                </a-layout-content>
            </a-layout>
        </a-layout>
    </a-layout>
</template>

<script lang="ts">
    import {Component, Vue, Watch} from 'vue-property-decorator'
    import {navList} from '@/router/manageRoutes'

    @Component
    export default class App extends Vue {
        private breadcrumbList: Array<string> = []
        private navList = navList

        @Watch('$route')
        private changeBreadCrumb() {
            this.breadcrumbList = this.$route.matched.filter(item => item.meta.alias).map(item => item.meta.alias)
        }

        private mounted() {
            this.changeBreadCrumb()
        }
    }
</script>

<style lang="sass" scoped>
#mainLayout
  height: 100vh
  .slider-head
    padding: 12px 12px 0 12px
  .wrap
    margin: 0
    padding: 14px
    min-height: 280px
    height: 100%
    overflow-y: scroll
    background: #fff
    &::-webkit-scrollbar
        width: 4px
    &::-webkit-scrollbar-thumb
        border-radius: 2px
        background-color: #e0e0e0
  .logo
    width: 140px
    height: 40px
    background: url(../../public/images/deepai.png) center no-repeat
    background-size: 140px 40px
</style>
