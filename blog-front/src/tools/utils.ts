/**
 * 拼接缩略图图片地址 （仅适用于七牛云）
 * @param imageUrl 后端返回的七牛图片地址
 * @param h 图片短边的长度 长边自适应
 */
export const sliceThumbnailImageUrl = (imageUrl: string, h: number = 100) => {
  if (imageUrl.startsWith('/b-oss/')) {
    return imageUrl + `?imageView2/0/h/${h}`
  }
  return imageUrl;
}