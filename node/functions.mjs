import sxes from 'sxes-compressor';

const {SxesGroup} = sxes;

async function load(uri) {
	const loadedGroup = new SxesGroup(uri);
	await loadedGroup.initialize();
	return loadedGroup;
}

async function getProjects(group) {
	const projects = group.projects;

	return Array.from(projects.keys()).map(key => {
		return projects.get(key).serialize();
	});
}

async function getAnalysis(group, uuid) {
	return group.getAnalysis(uuid);
}

async function getPosition(group, uuid) {
	const position = group.getPosition(uuid);
	await position.initialize()
	return position
}

async function getLine(group, uuid, typeName) {
	const position = await getPosition(group, uuid);
	return position.getType(typeName);
}

export default {
	load,
	getProjects,
	getAnalysis,
	getPosition,
	getLine
};