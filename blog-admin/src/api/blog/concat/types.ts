// 联系方式
export interface Concat {
  concatId: number;
  name: string;
  logoUrl: string;
  targetUrl: string;
  isMain: boolean,
  sort: number;
  status: number;
  createTime: string;
}

// 联系方式表单
export interface ConcatForm {
  concatId?: number;
  name: string;
  logoUrl: string;
  targetUrl: string;
  isMain: boolean;
  sort: number;
  status: number;
}

// 联系方式查询表单
export interface ConcatQueryForm {
  name?: string;
  createTimeBegin?: string;
  createTimeEnd?: string;
}