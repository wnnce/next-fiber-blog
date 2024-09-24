/// <reference types="vite/client" />

declare module 'nprogress';
declare module 'vditor/dist/method.min';

interface ImportMetaEnv {
  /* 后端请求地址 */
  readonly VITE_REQUEST_BASE_URL: string;
  /*七牛云Bucket的自定义域名*/
  readonly VITE_QINIU_DOMAIN: string;
  /* 远程资源需要被反向代理的前缀 */
  readonly VITE_REMOTE_SOURCE_PROXY_PREFIX: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}