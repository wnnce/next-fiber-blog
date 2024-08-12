import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "@/styles/globals.css";
import React from 'react'
import Header from '@/layouts/Header'
import Footer from '@/layouts/Footer'

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Create Next App",
  description: "Generated by create next app",
};

// 顶部导航列表
const navList = [
  { name: '博客', url: '#' },
  { name: '动态', url: '#' },
  { name: '分类', url: '/categorys' },
  { name: '标签', url: '/tags' },
  { name: '相册', url: '#' },
  { name: '关于', url: '#' },
  { name: '友情链接', url: '#' },
]

const RootLayout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  return (
    <html lang="zh-cn">
      <body className={inter.className}>
        <Header navList={navList} />
        <main>
          {children}
        </main>
        <Footer />
      </body>
    </html>
  );
}

export default RootLayout;