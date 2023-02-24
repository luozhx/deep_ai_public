import axios from 'axios'
import baseUrl from './baseUrl'

export function getJobsList() {
    return axios({
        url: `${baseUrl.deepAI}/job`,
        method: 'get'
    })
}

export function getJob(jobID: number) {
    return axios({
        url: `${baseUrl.deepAI}/job/${jobID}`,
        method: 'get'
    })
}

export function createNewJob(values: any) {
    return axios({
        url: `${baseUrl.deepAI}/job`,
        method: 'post',
        data: values
    })
}

export function deleteJob(jobID: number) {
    return axios({
        url: `${baseUrl.deepAI}/job/${jobID}`,
        method: 'delete'
    })
}

export interface IJobListItem {
    ID: number;
    Name: string;
    Framework: string;
    FrameworkVersion: string;
    // EntryPoint: string;
    // Args: string;
    Description: string;
    CreatedAt: string;
    Status: number;
    // Num: number;
    // CPU: number;
    // Memory: number;
    // GPU: number;
}
