import { observer } from "mobx-react-lite"
import { useEffect, useRef } from "react"
import { useNavigate } from "react-router"
import { useWebSocket } from "../../socket/Context"
import { sessionStore } from "../../store/SessionStore"

export default observer(() => {
	const navigate = useNavigate()
	const rootNode = useRef()

	return (
		<div ref={rootNode}>
			<h1>Available Rooms</h1>
			<button
				// onClick={() => {
				// 	navigate(`/room/${v4()}`)
				// }}
				onClick={() => {
					sessionStore.createRoom()
				}}
			>
				Create New Room
			</button>

			<ul>
				{sessionStore.rooms.map((room, index) => (
					<li key={index}>
						{room.uuid}
						<button onClick={() => {
							navigate(`/room/${room.uuid}`)
						}}>
							JOIN ROOM
						</button>
						<button onClick={() => sessionStore.leave()}>
							Leave
						</button>
						<ul>
							{room.clients.map((client, clientIndex) => (
								<li key={clientIndex} style={{marginLeft: "20px"}}>
									{client.uuid}
								</li>
							))}
						</ul>
					</li>
				))}
			</ul>
		</div>
	)
})
