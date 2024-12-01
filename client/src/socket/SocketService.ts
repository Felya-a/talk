import { makeAutoObservable } from "mobx"
import { sessionStore } from "../store/SessionStore"
import PingHandler from "./handlers/PingHandler"
import ShareRoomsHandler from "./handlers/ShareRoomsHandler"
import { MessagesTypes } from "./interface"

interface ReceiveMessage {
	type: MessagesTypes
	data: any
}

interface CreateRoomDto {
	room_name: string
}

interface JoinDto {
	room_uuid: string
}

type TransmitData = CreateRoomDto | JoinDto

interface TransmitMessage {
	type: MessagesTypes
	data: TransmitData
}



export interface MessageHandler {
	handle(data: any): Promise<void>
}

export enum SocketStatuses {
	CONNECTING = 0,
    OPEN = 1,
    CLOSING = 2,
    CLOSED = 3,
}

class SocketService {
	public socketStatus: SocketStatuses = SocketStatuses.CONNECTING
	private socket: WebSocket
	private handlers: Partial<Record<MessagesTypes, MessageHandler>>

	constructor(url: string) {
		makeAutoObservable(this)

		this.socket = new WebSocket(url)

		const pingHandler = new PingHandler()
		const shareRoomsHandler = new ShareRoomsHandler()

		this.handlers = {
			[MessagesTypes.PING]: pingHandler,
			[MessagesTypes.SHARE_ROOMS]: shareRoomsHandler
		}

		this.socket.onopen = this.onOpen.bind(this)
		this.socket.onclose = this.onClose.bind(this)
		this.socket.onmessage = this.onMessage.bind(this)
	}

	private async onOpen(event: Event): Promise<void> {
		console.log("onOpen", event)
		await new Promise(res => setTimeout(res, 1000)) // TODO ONLY DEBUG
		sessionStore.setSocket(this)
		this.socketStatus = SocketStatuses.OPEN
	}

	private async onClose(closeEvent: CloseEvent): Promise<void> {
		console.log("onClose", closeEvent)
		this.socketStatus = SocketStatuses.CLOSED
	}

	private async onMessage(message: MessageEvent): Promise<void> {
		console.log("Receive message: ", JSON.parse(message.data))
		const parsedMessage = JSON.parse(message.data) as ReceiveMessage

		const handler = this.handlers[parsedMessage.type]
		if (handler) {
			await handler.handle(parsedMessage.data)
		} else {
			console.warn(`No handler found for message type: ${parsedMessage.type}`)
		}
	}

	public send(message: TransmitMessage): void {
		this.socket.send(JSON.stringify(message))
	}

	public close(): void {
		this.socket.close()
	}

	// TODO: Возможно потом удалить
	// public getSocket(): WebSocket {
	// 	return this.socket
	// }
}

export default SocketService
