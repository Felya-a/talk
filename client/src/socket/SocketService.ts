import { makeAutoObservable } from "mobx"
import { sessionStore } from "../store/SessionStore"
import PingHandler from "./handlers/PingHandler"
import ShareRoomsHandler from "./handlers/ShareRoomsHandler"
import { MessagesTypes } from "./interface"

// TODO: Убрать отсюда интерфейсы
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

type LeaveDto = null

interface RelayIceDto {
	peer_id: string
	ice_candidate: string
}

interface RelaySdpDto {
	peer_id: string
	session_description: string
}

type MessageTypeToDtoMap = {
	[MessagesTypes.CREATE_ROOM]: CreateRoomDto
	[MessagesTypes.JOIN]: JoinDto
	[MessagesTypes.LEAVE]: LeaveDto
	[MessagesTypes.RELAY_ICE]: RelayIceDto
	[MessagesTypes.RELAY_SDP]: RelaySdpDto
}

type TransmitData<T extends MessagesTypes> = T extends keyof MessageTypeToDtoMap ? MessageTypeToDtoMap[T] : never

interface TransmitMessage<T extends MessagesTypes> {
	type: T
	data: TransmitData<T>
}

export interface MessageHandler {
	handle(data: any): Promise<void>
}

export enum SocketStatuses {
	CONNECTING = 0,
	OPEN = 1,
	CLOSING = 2,
	CLOSED = 3
}

class SocketService {
	public socketStatus: SocketStatuses = SocketStatuses.CONNECTING
	private socket: WebSocket
	private handlers: Partial<Record<MessagesTypes, MessageHandler>>
	private dynamicsHandlers: Partial<Record<MessagesTypes, (...args: any) => void>> = {}

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
		await new Promise(res => setTimeout(res, 500)) // TODO ONLY DEBUG
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

		// Динамические обработчики (подключенные после инициализации SocketService)
		const dynamicHandler = this.dynamicsHandlers[parsedMessage.type]
		if (dynamicHandler) {
			console.log("Найден динамический обработчик ", parsedMessage.type)
			dynamicHandler(parsedMessage.data)
			return
		}

		// Статические обработчики (подключенные в конструкторе)
		const handler = this.handlers[parsedMessage.type]
		if (handler) {
			await handler.handle(parsedMessage.data)
		} else {
			console.warn(`No handler found for message type: ${parsedMessage.type}`)
		}
	}

	public send<T extends MessagesTypes>(type: T, data: TransmitData<T>): void {
		console.log("Отправка сообщения ", type, data)
		const message: TransmitMessage<T> = { type, data }
		this.socket.send(JSON.stringify(message))
	}

	public on<T extends MessagesTypes>(type: T, handler: (...args: any) => void) {
		this.dynamicsHandlers[type] = handler
	}

	public close(): void {
		this.socket.close()
	}
}

export default SocketService
