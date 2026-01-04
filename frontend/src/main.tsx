import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'

let elementById = document.getElementById('root') ?? document.createElement('div');
ReactDOM.createRoot(elementById).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
