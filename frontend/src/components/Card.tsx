import './Card.css'

interface CardProps {
  children: React.ReactNode
  title?: string
  className?: string
  style?: React.CSSProperties
}

export default function Card({ children, title, className = '', style }: CardProps) {
  return (
    <div className={`card ${className}`} style={style}>
      {title && <h2 className="card-title">{title}</h2>}
      <div className="card-content">{children}</div>
    </div>
  )
}

