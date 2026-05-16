export type User={id:string;tenant_id:string;name:string;email:string;role:string};
export type Workflow={id:string;name:string;description:string;current_version:number;is_active:boolean;dag_definition?:DAGDefinition};
export type DAGDefinition={steps:Step[];edges?:Edge[];timeout_ms?:number};
export type Step={id:string;name:string;type:string;config?:Record<string,unknown>;depends_on?:string[]};
export type Edge={from:string;to:string;condition?:string};
export type Run={id:string;workflow_id:string;status:string;trigger_type:string;duration_ms:number;error_message:string};
export type RunStep={id:string;step_id:string;step_name:string;step_type:string;status:string;attempt_number:number;duration_ms:number;error_message:string;output?:unknown};
export type ApiResponse<T>={success:boolean;data:T;meta?:{page:number;limit:number;total:number;totalPages:number};error?:{code:string;message:string}};
