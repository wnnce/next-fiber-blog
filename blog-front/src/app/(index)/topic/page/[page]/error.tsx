'use client'

import React from 'react'
import { ErrorProps } from '@/app/layout'
import Empty from '@/components/Empty'
import DynamicCard from '@/components/DynamicCard'

const Error: React.FC<ErrorProps> = ({ error, reset }) => {
  return (
    <DynamicCard padding="1.5rem" title="TOOPICS" icon="i-tabler:world">
      <Empty text={error.message} icon="i-tabler:exclamation-circle" iconClassName="text-24 text-red-4" />
    </DynamicCard>
  )
}

export default Error;