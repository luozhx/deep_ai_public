import axios from 'axios'
import baseUrl from './baseUrl'

export function getDatasetList() {
    return axios({
        url: `${baseUrl.deepAI}/dataset`,
        method: 'get'
    })
}

export function getDataset(datasetID: number) {
    return axios({
        url: `${baseUrl.deepAI}/dataset/${datasetID}`,
        method: 'get'
    })
}

export function getDatasetDownloadUrl(datasetID: number) {
    return `${baseUrl.deepAI}/dataset/${datasetID}/download`
}

export function createNewDataset(values: FormData) {
    return axios({
        url: `${baseUrl.deepAI}/dataset`,
        method: 'post',
        data: values,
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    })
}

export function deleteDataset(datasetID: number) {
    return axios({
        url: `${baseUrl.deepAI}/dataset/${datasetID}`,
        method: 'delete'
    })
}

export interface IDatasetListItem {
    ID: number;
    Name: string;
    Size: number;
    Description: string;
    CreatedAt: string;
    Status: number;
    PVCStatus: number;
}
