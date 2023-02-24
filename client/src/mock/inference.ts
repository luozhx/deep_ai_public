import Mock from 'mockjs'
import baseUrl from '@/api/baseUrl'

function getInferenceList() {
    let dataList = []

    for (let i = 0; i < 30; i++) {
        dataList.push({
            'ID': Mock.Random.natural(),
            'CreatedAt': '@DATETIME("yyyy-MM-dd HH:mm:ss")',
            'UpdateAt': '@DATETIME("yyyy-MM-dd HH:mm:ss")',
            'DeleteAt': '@DATETIME("yyyy-MM-dd HH:mm:ss")',
            'Name': 'inference',
            'Description': 'inference',
            'Status|0-1': 1,
            'FrameWork': 'tf',
            'FrameWorkVersion': 'v1.0.1',
            'ModelPath': '/xxxx/xxxx',
            'Args': '-xxx',
            'Num': 0,
            'CPU': 1,
            'Memory': 100,
            'GPU': 1
        })
    }

    return dataList
}

Mock.mock(`${baseUrl.deepAI}/inference`, 'get', getInferenceList)
