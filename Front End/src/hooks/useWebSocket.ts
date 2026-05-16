import {useEffect} from 'react';
export function useWebSocket(onMessage:(event:any)=>void){useEffect(()=>{const token=localStorage.getItem('access_token');if(!token)return;const url=(import.meta.env.VITE_WS_URL||'ws://localhost:8080/ws')+`?token=${token}`;const ws=new WebSocket(url);ws.onmessage=(e)=>{try{onMessage(JSON.parse(e.data))}catch{}};return()=>ws.close();},[onMessage]);}
