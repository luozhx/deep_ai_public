import axios from 'axios'
import baseUrl from './baseUrl'

export function getInferenceList() {
    return axios({
        url: `${baseUrl.deepAI}/inference`,
        method: 'get'
    })
}

export function getInference(inferenceID: number) {
    return axios({
        url: `${baseUrl.deepAI}/inference/${inferenceID}`,
        method: 'get'
    })
}

export function createNewInference(values: any) {
    return axios({
        url: `${baseUrl.deepAI}/inference`,
        method: 'post',
        data: values
    })
}

export function deleteInference(inferenceID: number) {
    return axios({
        url: `${baseUrl.deepAI}/inference/${inferenceID}`,
        method: 'delete'
    })
}

export interface IInferenceListItem {
    ID: number;
    Name: string;
    Framework: string;
    FrameworkVersion: string;
    Description: string;
    CreatedAt: string;
    Status: number;
    ServiceName: string;
}
