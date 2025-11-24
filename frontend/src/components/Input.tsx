import './Input.css'

interface InputProps {
  type?: string
  placeholder?: string
  value: string
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void
  label?: string
  error?: string
  required?: boolean
}

export default function Input({
  type = 'text',
  placeholder,
  value,
  onChange,
  label,
  error,
  required = false,
}: InputProps) {
  return (
    <div className="input-group">
      {label && (
        <label className="input-label">
          {label}
          {required && <span className="required">*</span>}
        </label>
      )}
      <input
        type={type}
        placeholder={placeholder}
        value={value}
        onChange={onChange}
        className={`input ${error ? 'input-error' : ''}`}
        required={required}
      />
      {error && <span className="input-error-text">{error}</span>}
    </div>
  )
}


