/*****************************************************************************
 **                                                                         **
 **  I solved this problem in javascript first... so here's that solution.  **
 **                                                                         **
 **  To run it:                                                             **
 **      node /path/to/day13.js [ --example ]                               **
 **                                                                         **
 *****************************************************************************/

const fs = require('fs');

function part1(input) {
	let lines = input.split("\n");
	let pairs = [];
	for (let i = 0; i < lines.length; i+=3) {
		let a = JSON.parse(lines[i]);
		let b = JSON.parse(lines[i+1]);
		pairs.push([a, b])
	}

	let sum = 0;
	for (let i = 0; i < pairs.length; i++) {
		let res = compareLists(pairs[i][0], pairs[i][1]);
		if (res === -1){
			sum += i + 1;
		}
	}

	return sum;
}

function part2(input) {
	let lines = input.split("\n");
	let lists = [
		[[2]],
		[[6]],
	];
	for (let i = 0; i < lines.length; i++) {
		if (lines[i] !== "") {
			lists.push(JSON.parse(lines[i]));
		}
	}
	
	lists.sort(compareLists);
	
	let dividerPacketIndexes = [];
	for (let i = 0; i < lists.length; i++) {
		let v = lists[i];
		if (v.length == 1 && Array.isArray(v[0]) && v[0].length == 1 && (v[0][0] === 6 || v[0][0] === 2)) {
			dividerPacketIndexes.push(i+1);
			if (dividerPacketIndexes.length == 2) {
				break;
			}
		}
	}

	return dividerPacketIndexes[0] * dividerPacketIndexes[1];
}

function compareLists(a, b) {
	if (typeof a === 'number' && typeof b == 'number') {
		return a > b ? 1 : (a < b ? -1 : 0);
	} else if (Array.isArray(a) && Array.isArray(b)) {
		let minLen = Math.min(a.length, b.length)
		for (let i = 0; i < minLen; i++) {
			let cmp = compareLists(a[i], b[i]);
			if (cmp !== 0) {
				return cmp;
			}
		}
		return a.length < b.length ? -1 : (a.length > b.length ? 1 : 0)
	} else {
		if (typeof a === 'number') {
			return compareLists([a], b);
		} else {
			newA = a;
			newB = [b];
			return compareLists(a, [b]);
		}
	}
}


let inputFileName = 'day13.txt';

if (process.argv.includes('--example') || process.argv.includes('-e')) {
	inputFileName = 'day13_example1.txt';
}

const input = fs.readFileSync(`${__dirname}/${inputFileName}`).toString();

const part1Result = part1(input);
console.log('Part 1 result:', part1Result);

const part2Result = part2(input);
console.log('Part 2 result:', part2Result);
