import { makeAutoObservable } from "mobx"
import { MessagesTypes } from "../socket/interface"
import SocketService from "../socket/SocketService"

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
		this.rooms = rooms
	}

	createRoom() {
		this.socketService.send(MessagesTypes.CREATE_ROOM, {
			room_name: "test_from_react"
		})
	}

	leave() {
		this.socketService.send(MessagesTypes.LEAVE, null)
	}
}

export const sessionStore = new SessionStore();