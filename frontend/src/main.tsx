import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'

// 抑制Wails内部的dispatchWailsEvent错误
if (typeof window !== 'undefined') {
  // 拦截所有TypeError
  const originalError = console.error
  console.error = (...args) => {
    // 将参数转换为字符串检查
    const errorMsg = args.map(arg => {
      if (typeof arg === 'string') return arg
      if (arg instanceof Error) return arg.message
      try {
        return JSON.stringify(arg)
      } catch {
        return String(arg)
      }
    }).join(' ')

    // 过滤Wails相关错误
    if (errorMsg.includes('dispatchWailsEvent') ||
        errorMsg.includes('_wails.') ||
        errorMsg.includes('mac:WindowDidUpdate')) {
      return // 静默忽略
    }

    originalError.apply(console, args)
  }
}

let elementById = document.getElementById('root') ?? document.createElement('div');
ReactDOM.createRoot(elementById).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
