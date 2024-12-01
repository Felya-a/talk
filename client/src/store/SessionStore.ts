import { makeAutoObservable } from "mobx"
import { MessagesTypes } from "../socket/interface"
import SocketService, { SocketStatuses } from "../socket/SocketService"

interface Client {
	uuid: string
}

interface Room {
	clients: Client[]
	uuid: string
	name: string
}

class SessionStore {
	socketService: SocketService
	rooms: Room[] = []

	constructor() {
		makeAutoObservable(this)
	}

	setSocket(socketService: SocketService) {
		this.socketService = socketService
	}

	updateRooms(rooms: Room[]) {
		console.log("updateRooms: ", rooms)
		this.rooms = rooms
	}

	createRoom() {
		this.socketService.send({
			type: MessagesTypes.CREATE_ROOM,
			data: {
				room_name: "test_from_react"
			}
		})
	}

	joinToRoom(roomUuid) {
		this.socketService.send({
			type: MessagesTypes.JOIN,
			data: {
				room_uuid: roomUuid
			}
		})
	}

	leave() {
		this.socketService.send({
			type: MessagesTypes.LEAVE,
			data: null
		})
	}
}

export const sessionStore = new SessionStore();