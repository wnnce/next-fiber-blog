import type { TableRowSelection } from '@arco-design/web-vue'

interface ConstantPool{
  readonly pageSizeOption: number[]
}

// 导出的常量
export const constant: ConstantPool = {
  pageSizeOption: [10, 20, 40]
}

// token本地保存的key
export const TOKEN_KEY: string = 'Authorization_Bearer_Token';

// 登录用户信息本地保存的key
export const LOCAl_USER_KEY: string = 'Local_Login_User';

// 表格选择行参数
export const ROW_SELECTION: TableRowSelection = {
  type: 'checkbox',
  showCheckedAll: true,
  onlyCurrent: false
}