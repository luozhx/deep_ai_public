import axios from 'axios'
import baseUrl from './baseUrl'

export function getModelList() {
    return axios({
        url: `${baseUrl.deepAI}/model`,
        method: 'get'
    })
}

export function getModel(modelID: number) {
    return axios({
        url: `${baseUrl.deepAI}/model/${modelID}`,
        method: 'get'
    })
}

export function getModelFileList(modelID: number, path: string) {
    return axios({
        url: `${baseUrl.deepAI}/model/${modelID}/dir`,
        method: 'get',
        params: {
            path
        }
    })
}

export function getModelDownloadUrl(modelID: number) {
    return `${baseUrl.deepAI}/model/${modelID}/download`
}

export function createNewModel(values: any) {
    return axios({
        url: `${baseUrl.deepAI}/model`,
        method: 'post',
        data: values
    })
}

export function deleteModel(modelID: number) {
    return axios({
        url: `${baseUrl.deepAI}/model/${modelID}`,
        method: 'delete'
    })
}

export interface IModelFileListItem {
    name: string;
    isDir: boolean;
}

export interface IModelListItem {
    ID: number;
    Name: string;
    Framework: string;
    FrameworkVersion: string;
    Description: string;
    CreatedAt: string;
    Status: number;
    PVCStatus: number;
}
