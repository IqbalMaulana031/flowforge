import {create} from 'zustand';import type {User} from '../types';
type AuthState={user?:User;setSession:(s:{access_token:string;refresh_token:string;user:User})=>void;logout:()=>void};
export const useAuth=create<AuthState>((set)=>({user:undefined,setSession:(s)=>{localStorage.setItem('access_token',s.access_token);localStorage.setItem('refresh_token',s.refresh_token);set({user:s.user});},logout:()=>{localStorage.removeItem('access_token');localStorage.removeItem('refresh_token');set({user:undefined});}}));
