import {useEffect,useState} from 'react';import {health} from '../../api/health.api';
export function HealthPanel(){const [data,setData]=useState<unknown>();useEffect(()=>{health().then(setData).catch(e=>setData(e.message))},[]);return <section className="card"><h2>Health</h2><pre>{JSON.stringify(data,null,2)}</pre></section>}
