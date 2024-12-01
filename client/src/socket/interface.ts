export enum MessagesTypes {
	PING = "ping",
	PONG = "pong",
	CREATE_ROOM = "create-room",
	JOIN = "join",
	LEAVE = "leave",
	SHARE_ROOMS = "share-rooms",
	ADD_PEER = "add-peer",
	REMOVE_PEER = "remove-peer",
	RELAY_SDP = "relay-sdp",
	RELAY_ICE = "relay-ice",
	ICE_CANDIDATE = "ice-candidate",
	SESSION_DESCRIPTION = "session-description"
}

export interface MessageJoin {
	type: MessagesTypes.JOIN
	data: {
		room_uuid: string
	}
}

export interface MessageCreateRoom {
	type: MessagesTypes.CREATE_ROOM
	data: {
		room_name: string
	}
}

export interface MessageLeave {
	type: MessagesTypes.LEAVE
	data: null
}

export interface MessageRelaySdp {
	type: MessagesTypes.RELAY_SDP
	data: {
        peer_id: string
        session_description: string
    }
}

export interface MessageRelayIce {
	type: MessagesTypes.RELAY_ICE
	data: {
        peer_id: string
        ice_candidate: string
    }
}

export interface Message {
	type: MessagesTypes
	data: any
}
