'use client'

import React from 'react'
import { AppProgressBar as ProgressBar } from 'next-nprogress-bar'

const GlobalProgressBar: React.FC = () => {
  return (
    <ProgressBar
      height="3px"
      color="#46B952"
      options={{ showSpinner: false }}
      shallowRouting
    />
  )
}

export default GlobalProgressBar;