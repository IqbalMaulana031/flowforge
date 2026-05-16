import {api} from './client';import type {ApiResponse,Workflow,DAGDefinition} from '../types';
export async function listWorkflows(){const {data}=await api.get<ApiResponse<Workflow[]>>('/workflows');return data.data}
export async function getWorkflow(id:string){const {data}=await api.get<ApiResponse<Workflow>>(`/workflows/${id}`);return data.data}
export async function createWorkflow(payload:{name:string;description:string;dag_definition:DAGDefinition}){const {data}=await api.post<ApiResponse<Workflow>>('/workflows',payload);return data.data}
export async function updateWorkflow(id:string,payload:{name:string;description:string;dag_definition:DAGDefinition}){const {data}=await api.put<ApiResponse<Workflow>>(`/workflows/${id}`,payload);return data.data}
