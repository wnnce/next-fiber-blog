import '@/styles/components/links.scss'
import React from 'react'
import DynamicCard from '@/components/DynamicCard'
import { listLinks } from '@/lib/api'
import Empty from '@/components/Empty'
import RichImage from '@/components/RichImage'

const Links: React.FC = async () => {
  const { code, message, data: links } = await listLinks();

  if (code !== 200) {
    throw new Error(message);
  }
  return (
    <DynamicCard padding="1.5rem" title={`LINKS${links && links.length > 0 ? ` (${links.length})` : ''}`}
                 icon="i-tabler:link"
                 multiple={40}
    >
      { !links || links.length == 0 ? (
        <Empty text="还没有友情链接..." iconClassName="text-24" />
      ) : (
        <section>
          <ul className="list-none flex flex-wrap gap-4 link-list-ul">
            { links.map(link => (
              <li key={link.linkId} className="p-3 flex block flex-col gap-row-2 link-list-li">
                <div className="flex justify-between gap-col-2 items-center">
                  <h3 className="font-bold a-hover-line-text-md"><a href={link.targetUrl} target="_blank">{ link.name }</a></h3>
                  <RichImage className="link-avatar" src={link.coverUrl} width={40} height={40} fill radius="50%" thumbnail />
                </div>
                <p className="text-sm info-text line-clamp-2">{link.summary}</p>
              </li>
            ))}
          </ul>
          <h2 className="mt-8 font-bold pl-2 border-l-4 link-register-title">
            申请格式
          </h2>
          <p className="info-text mt-4">[ 博客名称 + 博客地址 + Logo地址 + 博客简介 ]</p>
          <p className="desc-text text-xs mt-2">请在下方留言板留下你的信息...</p>
        </section>
      )}
    </DynamicCard>
  )
}

export default Links;