import '@/styles/components/button.scss';
import React, { CSSProperties } from 'react'

interface Props {
  icon?: string;
  text?: string;
  loading?: boolean;
  disabled?: boolean;
  onClick?: () => void;
  className?: string;
  style?: CSSProperties;
  type?: 'submit' | 'reset' | 'button'
}

const Button: React.FC<Props> = ({
  icon,
  text = '按钮',
  loading = false,
  disabled = false,
  onClick,
  className,
  style,
  type
}) => {

  (!disabled && loading) && (disabled = true);

  let buttonIcon: React.ReactNode

  if (loading) {
    buttonIcon = <i className="inline-block text-3.5 i-tabler:loader-2 animate-spin" />
  } else if (icon) {
    buttonIcon = <i className={`inline-block text-3.5 ${icon}`} />
  }
  return (
    <button style={style}
            className={`button-component px-3 text-xs ${disabled ? 'cursor-not-allowed button-component-disabled ' : 'button-component-active'} ${className ? className : ''}`}
            disabled={disabled}
            onClick={onClick}
            type={type}
    >
      <div className="flex gap-col-1 items-center">
        { buttonIcon }
        {text}
      </div>
    </button>
  )
}


export default Button;