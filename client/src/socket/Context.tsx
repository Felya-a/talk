// WebSocketContext.tsx
import React, { createContext, useContext, useEffect, useRef } from "react";
import SocketService, { SocketStatuses } from "./SocketService";
import { sessionStore } from "../store/SessionStore"
import { observe } from "mobx"
import { observer } from "mobx-react-lite"

const WebSocketContext = createContext<SocketService | null>(null);

interface WebSocketProviderProps {
    children: React.ReactNode;
    url: string; // WebSocket URL
}

export const WebSocketProvider: React.FC<WebSocketProviderProps> = observer(({ children, url }) => {
    const webSocketServiceRef = useRef<SocketService | null>(null);

    useEffect(() => {
        // Инициализируем SocketService
        webSocketServiceRef.current = new SocketService(url);

        return () => {
            // Закрываем соединение при размонтировании
            webSocketServiceRef.current?.close();
        };
    }, [url]);

    if (webSocketServiceRef?.current?.socketStatus === SocketStatuses.CONNECTING) {
        return (
            <div>
                Подключение к серверу...
            </div>
        )
    }

    if (webSocketServiceRef?.current?.socketStatus === SocketStatuses.CLOSED) {
        return (
            <div>
                Связь с сервером потеряна
            </div>
        )
    }

    return (
        <WebSocketContext.Provider value={webSocketServiceRef.current}>
            {children}
        </WebSocketContext.Provider>
    );
});

export const useWebSocket = (): SocketService | null => {
    return useContext(WebSocketContext);
};
