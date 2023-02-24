import RouterView from '@/components/RouterView.vue'

import Console from '@/views/Console.vue'

import CodeList from '@/views/Code/CodeList.vue'
import NewCode from '@/views/Code/NewCode.vue'

import DatasetList from '@/views/Dataset/DatasetList.vue'
import NewDataset from '@/views/Dataset/NewDataset.vue'

import ModelList from '@/views/Model/ModelList.vue'
import NewModel from '@/views/Model/NewModel.vue'

import JobList from '@/views/Job/JobList.vue'
import NewJob from '@/views/Job/NewJob.vue'

import InferenceList from '@/views/Inference/InferenceList.vue'
import NewInference from '@/views/Inference/NewInference.vue'

import DemoView from '@/views/Demo/DemoView.vue'
import DemoMnist from '@/views/Demo/DemoMnist.vue'

import ImageLabel from '@/views/ImageLabel/ImageLabel.vue'
import DemoEDSR from '@/views/Demo/DemoEDSR.vue'

const manageRoutes = [{
    path: 'console',
    name: 'Console',
    component: Console,
    meta: {
        alias: '状态监控',
        icon: 'desktop'
    }
}, {
    path: 'code',
    component: RouterView,
    redirect: {
        name: 'codeList'
    },
    children: [{
        path: 'list',
        name: 'codeList',
        component: CodeList
    }, {
        path: 'new',
        name: 'newCode',
        component: NewCode
    }],
    meta: {
        alias: '代码管理',
        icon: 'tool'
    }
}, {
    path: 'dataset',
    component: RouterView,
    redirect: {
        name: 'datasetList'
    },
    children: [{
        path: 'list',
        name: 'datasetList',
        component: DatasetList
    }, {
        path: 'new',
        name: 'newDataset',
        component: NewDataset
    }],
    meta: {
        alias: '数据集管理',
        icon: 'hdd'
    }
}, {
    path: 'model',
    component: RouterView,
    redirect: {
        name: 'modelList'
    },
    children: [{
        path: 'list',
        name: 'modelList',
        component: ModelList
    }, {
        path: 'new',
        name: 'newModel',
        component: NewModel
    }],
    meta: {
        alias: '模型管理',
        icon: 'code'
    }
}, {
    path: 'job',
    component: RouterView,
    redirect: {
        name: 'jobList'
    },
    children: [{
        path: 'list',
        name: 'jobList',
        component: JobList
    }, {
        path: 'new',
        name: 'newJob',
        component: NewJob
    }],
    meta: {
        alias: '作业管理',
        icon: 'code'
    }
}, {
    path: 'inference',
    component: RouterView,
    redirect: {
        name: 'inferenceList'
    },
    children: [{
        path: 'list',
        name: 'inferenceList',
        component: InferenceList
    }, {
        path: 'new',
        name: 'newInference',
        component: NewInference
    }],
    meta: {
        alias: '推理服务',
        icon: 'rocket'
    }
}, {
    path: 'demo',
    component: RouterView,
    redirect: {
        name: 'demoView'
    },
    children: [{
        path: 'demoView',
        name: 'demoView',
        component: DemoView
    }, {
        path: 'mnist',
        name: 'demoMnist',
        component: DemoMnist,
        meta: {
            alias: 'mnist'
        }
    }, {
        path: 'edsr',
        name: 'demoEDSR',
        component: DemoEDSR,
        meta: {
            alias: 'EDSR'
        }
    }],
    meta: {
        alias: '示范案例',
        icon: 'codepen'
    }
}, {
    path: 'imageLabel',
    component: RouterView,
    redirect: {
        name: 'imageLabel'
    },
    children: [{
        path: 'imageLabel',
        name: 'imageLabel',
        component: ImageLabel
    }],
    meta: {
        alias: '图像标注',
        icon: 'codepen'
    }
}]

export const navList = manageRoutes.map(item => ({
    path: item.path,
    name: item.meta.alias,
    icon: item.meta.icon
}))

export default manageRoutes
