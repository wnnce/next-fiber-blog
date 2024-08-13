/** @type {import('next').NextConfig} */
const nextConfig = {
  images: {
    remotePatterns: [
      { hostname: 'file.qiniu.vnc.ink' }
    ]
  }
};

export default nextConfig;
