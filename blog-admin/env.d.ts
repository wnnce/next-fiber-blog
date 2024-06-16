/// <reference types="vite/client" />

interface ImportMetaEnv {
  /* 后端请求地址 */
  readonly VITE_REQUEST_BASE_URL: string;
  /*七牛云Bucket的自定义域名*/
  readonly VITE_QINIU_DOMAIN: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}