import '@/styles/components/timeline.scss'
import React from 'react'
import { formatDateTime } from '@/tools/utils'

export const Timeline: React.FC<{
  children: React.ReactNode
}> = ({ children }) => {
  return (
    <ul className="list-none timeline">
      { children }
    </ul>
  )
}

export const TimeLineItem: React.FC<{
  time: Date | string
  children: React.ReactNode
}> = ({ time, children }) => {
  return (
    <li className="timeline-item">
      <div className="timeline-item-time desc-text line-height-loose flex items-center py-2">
        <i className="inline-block i-tabler:clock mr-1.5"></i>
        <time dateTime={time.toString()} className="text-xs ">{ formatDateTime(time.toString()) }</time>
      </div>
      <div className="p-4 timeline-item-body">
        { children }
      </div>
    </li>
  )
}