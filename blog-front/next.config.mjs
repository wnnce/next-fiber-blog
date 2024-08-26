/** @type {import('next').NextConfig} */
const nextConfig = {
  rewrites() {
    return [
      // 首页重定向到第一页
      { source: '/', destination: '/page/1' },
      // 重定向到动态第一页
      { source: '/topic', destination: '/topic/page/1' }
    ]
  },
  images: {
    remotePatterns: [
      { hostname: 'file.qiniu.vnc.ink' }
    ]
  },
};

export default nextConfig;
