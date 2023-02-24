import axios from 'axios'
import baseUrl from '@/api/baseUrl'

export function getSystemStatus(params: any = {}) {
    return axios({
        url: `${baseUrl.deepAI}/system/status`,
        method: 'get',
        params
    })
}
