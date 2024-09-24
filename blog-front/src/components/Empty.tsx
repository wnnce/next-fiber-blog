import React from 'react'

const Empty: React.FC<{
  text?: string;
  icon?: string;
  iconSize?: string;
  iconClassName?: string;
  textSize?: string;
  textClassName?: string
}> = ({
  text = '暂无数据',
  icon = 'i-tabler:template-off',
  iconSize,
  iconClassName,
  textSize,
  textClassName
}) => {
  return (
    <div className="text-center info-text">
      <i className={`${iconClassName ? iconClassName : ''} inline-block ${icon}`}
         style={ iconSize ? { fontSize: iconSize } : undefined }
      />
      <p className={`${textClassName ? textClassName : ''} text-center mt-4`}
         style={ textSize ? { fontSize: textSize } : undefined }
      >
        { text }
      </p>
    </div>
  )
}

export default Empty;