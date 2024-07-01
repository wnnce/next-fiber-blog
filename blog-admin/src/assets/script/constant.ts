interface ConstantPool{
  pageSizeOption: readonly number[]
}

// 导出的常量
export const constant: ConstantPool = {
  pageSizeOption: [10, 20, 40]
}