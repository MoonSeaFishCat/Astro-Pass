import './Button.css'

interface ButtonProps {
  children: React.ReactNode
  onClick?: () => void
  type?: 'button' | 'submit' | 'reset'
  variant?: 'primary' | 'secondary' | 'outline'
  disabled?: boolean
  fullWidth?: boolean
  style?: React.CSSProperties
  size?: 'small' | 'medium' | 'large'
}

export default function Button({
  children,
  onClick,
  type = 'button',
  variant = 'primary',
  disabled = false,
  fullWidth = false,
  style,
  size = 'medium',
}: ButtonProps) {
  return (
    <button
      type={type}
      onClick={onClick}
      disabled={disabled}
      style={style}
      className={`btn btn-${variant} ${fullWidth ? 'btn-full' : ''} btn-${size}`}
    >
      {children}
    </button>
  )
}

