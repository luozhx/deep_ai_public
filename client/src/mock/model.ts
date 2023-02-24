import Mock from 'mockjs'
import baseUrl from '@/api/baseUrl'

function getModelList() {
    let dataList = []

    for (let i = 0; i < 30; i++) {
        dataList.push({
            'ID': Mock.Random.natural(),
            'CreatedAt': '@DATETIME("yyyy-MM-dd HH:mm:ss")',
            'UpdateAt': '@DATETIME("yyyy-MM-dd HH:mm:ss")',
            'DeleteAt': '@DATETIME("yyyy-MM-dd HH:mm:ss")',
            'Name': 'test',
            'Description': 'test',
            'Status|0-1': 0,
            'FrameWork': 'tf',
            'FrameWorkVersion': 'v1.0.1',
            'CodePath': '/xxxx/xxxx',
            'EntryPoint': '/xxxx',
            'OutputPath': '/xxxx',
            'Args': '-xxx',
            'DatasetPath': '/xxx',
            'Num': 0,
            'CPU': 1,
            'Memory': 100,
            'GPU': 1,
            'TrainTime': '02:27:37.09'
        })
    }

    return dataList
}

Mock.mock(`${baseUrl.deepAI}/model`, 'get', getModelList)
