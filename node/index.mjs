import Websocket from 'ws';
import bakaRpc from 'baka-rpc';
import functions from './functions.mjs';

const {BakaRPC, constants} = bakaRpc;

const ws = new Websocket(`ws://localhost:${process.argv[2] !== undefined ? process.argv[2] : 9999}`);
const duplex = Websocket.createWebSocketStream(ws, {encoding: 'utf8'});
const rpc = new BakaRPC(duplex, duplex, {});

let loadedGroup = undefined;

rpc.on('load', async (uri) => {
	loadedGroup = await functions.load(uri);
	return true;
});

rpc.on('getProjects', async () => {
	return await functions.getProjects(loadedGroup);
});

rpc.on('getAnalysis', async (uuid) => {
	return (await functions.getAnalysis(loadedGroup, uuid)).serialize();
});

rpc.on('getPosition', async (uuid) => {
	return (await functions.getPosition(loadedGroup, uuid)).serialize();
});

rpc.on('getLine', async (uuid, typeName) => {
	return (await functions.getLine(loadedGroup, uuid, typeName)).serialize();
});

/*
functions.load('E:/probelab/work/work.plzip').then(async (group) => {
	const test = await functions.getPosition(group, '943940a8-17b1-43ce-b3b1-b49bfcadc9c1');

	console.log();
})
*/