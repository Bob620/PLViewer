import assert from 'assert';
import mocha from 'mocha';

import functions from '../functions.mjs';

const {describe, it, before} = mocha;

describe('Functions', async () => {
	let group;

	before(async () => {
		group = await functions.load('E:/probelab/work/work.plzip');
	});

	describe('#load', async () => {
		it('should load and initialize a plzip', async () => {
			await assert.doesNotReject(functions.load.bind(undefined, 'E:/probelab/work/work.plzip'));
		});
	});

	describe('#getProjects', async () => {

		it('???', async () => {
			await assert.doesNotReject(functions.getProjects.bind(undefined, group));
		});
	});

	describe('#getAnalysis', async () => {
		it('???', async () => {
			await assert.doesNotReject(functions.getAnalysis.bind(undefined, group, '0b6fe181-d096-477f-9ea9-d1ab89795c28'));
		});
	});

	describe('#getPosition', async () => {
		it('???', async () => {
			const position = await functions.getPosition(group, '00a9f573-f461-4a91-8148-422f806abf9d');

			await assert.deepStrictEqual(await position.serialize(), {
				uuid: '00a9f573-f461-4a91-8148-422f806abf9d',
				comment: 'MAC_TeO2_10kV100nA_15m_30s_30x : Pos. : 15',
				operator: 'AVDH',
				background: '',
				condition: 'f60e6ebfcd2b8c8c33c401323626d6194d33e80c7b0e725c4ee3666156307fda',
				rawCondition: 'dbe8b67bb247e96d06f161bfa9d8ead4bb9cf9725ec4cf69d7ae9943a879645c',
				types: [
					'qlw'
				]
			});
		});
	});

	describe('#getLine', async () => {
		it('???', async () => {
			await assert.doesNotReject(functions.getLine.bind(undefined, group, '00a9f573-f461-4a91-8148-422f806abf9d', 'qlw'));
		});
	});
});