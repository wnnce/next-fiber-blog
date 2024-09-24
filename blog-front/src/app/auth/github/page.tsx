'use client'

import React, { useEffect, useRef, useState } from 'react'
import { clientAuthTokenKey, loginWithGithub } from '@/lib/client-api'
import { useRouter } from 'next/navigation';

const Page: React.FC = () => {
  const router = useRouter();

  const [message, setMessage] = useState<string>('');
  const init = useRef<boolean>(false);
  const login = async (code: string, state: string) => {
    const { code: status, message, data } = await loginWithGithub(code);
    if (status != 200) {
      setMessage(message);
      return;
    } else {
      localStorage.setItem(clientAuthTokenKey, data);
      const path = atob(state);
      router.push(path);
    }
  }

  useEffect(() => {
    if (!init.current) {
      init.current = true;
      const params = new URLSearchParams(window.location.search);
      const code = params.get('code');
      const state = params.get('state');
      if (!code || !state) {
        setMessage('登录参数错误');
      } else {
        login(code, state);
      }
    }
  }, [])

  return (
    <div className="fixed top-0 left-0 right-0 bottom-0 p-4 z-20 bg-gray-2 text-black text-center">
      <i className="inline-block text-20 i-tabler:brand-github" />
      { message ? (
        <p className="text-red-5">
          {message}
        </p>
      ) : (
        <>
          <p className="line-height-loose text-lg">
            Github登录中...
          </p>
          <i className="inline-block text-6 i-tabler:loader-2 animate-spin" />
        </>
      )}
    </div>
  )
}

export default Page;