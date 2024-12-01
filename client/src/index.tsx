import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';
import { WebSocketProvider } from './socket/Context'

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);
root.render(
  <React.StrictMode>
    <WebSocketProvider url="ws://localhost:8090/ws">
      <App />
    </WebSocketProvider>
  </React.StrictMode>
);
