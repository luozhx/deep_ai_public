import Mock from 'mockjs'
import baseUrl from '@/api/baseUrl'

function getCodeList() {
    let dataList = []

    for (let i = 0; i < 30; i++) {
        dataList.push({
            'ID': Mock.Random.natural(),
            'CreatedAt': '@DATETIME("yyyy-MM-dd HH:mm:ss")',
            'UpdateAt': '@DATETIME("yyyy-MM-dd HH:mm:ss")',
            'DeleteAt': '@DATETIME("yyyy-MM-dd HH:mm:ss")',
            'Name': 'code',
            'Description': 'code',
            'Status|0-1': 0,
            'FrameWork': 'tf',
            'FrameWorkVersion': 'v1.0.1',
            'StoragePath': '/xxxx/xxxx'
        })
    }

    return dataList
}

Mock.mock(`${baseUrl.deepAI}/code`, 'get', getCodeList)
