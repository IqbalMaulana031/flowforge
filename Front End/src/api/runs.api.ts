import {api} from './client';import type {ApiResponse,Run,RunStep} from '../types';
export async function triggerRun(workflow_id:string){const {data}=await api.post<ApiResponse<Run>>('/runs/trigger',{workflow_id,payload:{source:'dashboard'}});return data.data}
export async function listRuns(){const {data}=await api.get<ApiResponse<Run[]>>('/runs');return data.data}
export async function getRunSteps(id:string){const {data}=await api.get<ApiResponse<RunStep[]>>(`/runs/${id}/steps`);return data.data}
export async function getRunLogs(id:string){const {data}=await api.get<ApiResponse<unknown[]>>(`/runs/${id}/logs`);return data.data}
export async function cancelRun(id:string){await api.delete(`/runs/${id}`)}
