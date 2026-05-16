import {api} from './client';
import {encryptPassword} from './password-crypto';
import type {ApiResponse, User} from '../types';

type AuthPayload = {access_token: string; refresh_token: string; user: User};

export async function register(payload: {tenant_name: string; tenant_slug: string; name: string; email: string; password: string}) {
  const encryptedPassword = await encryptPassword(payload.password);
  const {data} = await api.post<ApiResponse<AuthPayload>>('/auth/register', {
    tenant_name: payload.tenant_name,
    tenant_slug: payload.tenant_slug,
    name: payload.name,
    email: payload.email,
    encrypted_password: encryptedPassword,
  });
  return data.data;
}

export async function login(payload: {email: string; password: string}) {
  const encryptedPassword = await encryptPassword(payload.password);
  const {data} = await api.post<ApiResponse<AuthPayload>>('/auth/login', {
    email: payload.email,
    encrypted_password: encryptedPassword,
  });
  return data.data;
}

export async function profile() {
  const {data} = await api.get<ApiResponse<User>>('/auth/profile');
  return data.data;
}
