import './Loading.css'

interface LoadingProps {
  size?: 'small' | 'medium' | 'large'
  text?: string
}

export default function Loading({ size = 'medium', text }: LoadingProps) {
  return (
    <div className="loading-container">
      <div className={`loading-spinner loading-${size}`}>
        <div className="spinner-circle"></div>
      </div>
      {text && <p className="loading-text">{text}</p>}
    </div>
  )
}


