import { sendDelete, sendGet, sendPost, sendPut } from '@/api/request'
import type { Category, CategoryForm, CategoryUpdateForm } from '@/api/blog/category/types'

export const categoryApi = {
  /**
   * 获取分类树型结构数据
   */
  manageTree: () => {
    return sendGet<Category[]>('/category/manage/tree');
  },
  /**
   * 保存分类
   * @param form 分类参数
   */
  saveCategory: (form: CategoryForm) => {
    return sendPost<null>('/category', form)
  },
  /**
   * 更新分类
   * @param form
   */
  updateCategory: (form: CategoryForm) => {
    return sendPut<null>('/category', form)
  },
  /**
   * 快捷更新分类
   * @param form 快捷更新分类表单
   */
  updateCategoryStatus: (form: CategoryUpdateForm) => {
    return sendPut<null>('/category/status', form)
  },
  /**
   * 删除分类
   * @param id 分类Id
   */
  deleteCategory: (id: number) => {
    return sendDelete<null>(`/category/${id}`)
  },
  tree: () => {
    return sendGet<Category[]>('/open/category/list')
  }
}