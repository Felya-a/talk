import { MessageHandler } from "../SocketService"
import { sessionStore } from "../../store/SessionStore"

export default class ShareRoomsHandler implements MessageHandler {
	async handle(data: any) {
        sessionStore.updateRooms(data)
    }
}
