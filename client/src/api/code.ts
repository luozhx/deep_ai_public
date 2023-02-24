import axios from 'axios'
import baseUrl from './baseUrl'

export function getCodeList(params: any = {}) {
    return axios({
        url: `${baseUrl.deepAI}/code`,
        method: 'get',
        params
    })
}

export function getCodeFileList(codeID: number, path: string) {
    return axios({
        url: `${baseUrl.deepAI}/code/${codeID}/dir`,
        method: 'get',
        params: {
            path
        }
    })
}

export function getCode(codeID: number) {
    return axios({
        url: `${baseUrl.deepAI}/code/${codeID}`,
        method: 'get'
    })
}

export function getCodeNotebookUrl(codeID: number, serviceName: string) {
    return `${baseUrl.deepAI}/code/${codeID}/notebook/${serviceName}/`
}

export function getCodeDownloadUrl(codeID: number) {
    return `${baseUrl.deepAI}/code/${codeID}/download`
}

export function createNewCode(values: any) {
    return axios({
        url: `${baseUrl.deepAI}/code`,
        method: 'post',
        data: values
    })
}

export function deleteCode(codeID: number) {
    return axios({
        url: `${baseUrl.deepAI}/code/${codeID}`,
        method: 'delete'
    })
}

export interface ICodeFileListItem {
    name: string;
    isDir: boolean;
}

export interface ICodeListItem {
    ID: number;
    Name: string;
    Description: string;
    ServiceName: string;
    Framework: string;
    FrameworkVersion: string;
    CreatedAt: string;
    Status: number;
    PVCStatus: number;
    NotebookStatus: number;
}
