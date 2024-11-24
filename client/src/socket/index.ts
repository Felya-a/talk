import {io, ManagerOptions} from 'socket.io-client';

const options: Partial<ManagerOptions> = {
  forceNew: true,
  reconnection: true,
  timeout : 10000, // before connect_error and connect_timeout are emitted.
  transports : ["websocket"]
}

const socket = io('ws://192.168.0.2:3001', options);

export default socket;