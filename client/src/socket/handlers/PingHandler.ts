import { MessageHandler } from "../SocketService"

export default class PingHandler implements MessageHandler {
	async handle(data: any) {
		console.log(data)
	}
}
