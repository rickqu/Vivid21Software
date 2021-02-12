let width = 1000;
let height = 1000;
let paused = false;

let canvas = d3.select('body').append('svg')
	.attr('width', width)
	.attr('height', height)
	.append('g');


canvas.append('circle')
	.attr('r', 2)
	.style('fill', 'rgba(0, 0, 0, 1.0)')
	.attr('transform', 'translate(' + width / 2 + ', ' + height / 2 + ')');

for (let i = 1; i < 15; i++) {
	canvas.append('circle')
		.attr('r', 70 * i)
		.style('fill', 'rgba(0, 0, 0, 0)')
		.style('border', 'rgba(0, 0, 0, 1.0)')
		.style('fill', 'rgba(0, 0, 0, 0)')
		.style('stroke', 'rgba(0, 0, 0, 0.2)')
		.attr('transform', 'translate(' + width / 2 + ', ' + height / 2 + ')');
}

canvas.append('rect')
	.attr('x', 1)
	.attr('y', 1)
	.attr('width', width - 2)
	.attr('height', height - 2)
	.style('border', 'rgba(0, 0, 0, 1.0)')
	.style('fill', 'rgba(0, 0, 0, 0)')
	.style('stroke', 'rgba(0, 0, 0, 1.0)')
	.style('stroke-width', '1');

const socket = new WebSocket('ws://127.0.0.1:9000/ws');
const numChainsPerFern = 8;
const distBetweenLEDs = 5;

// var exdata = {
// 	ferns: [
// 		{
// 			location: {x: 50, y: 50},
// 			leds: [
// 				[{R: 100, G: 100, B: 100}, {R: 100, G: 100, B: 100}, {R: 100, G: 150, B: 100}, {R: 100, G: 100, B: 100}],
// 				[{R: 200, G: 200, B: 200}, {R: 200, G: 200, B: 200}, {R: 200, G: 200, B: 200}, {R: 200, G: 200, B: 200}],
// 			]
// 		},
// 	],
// 	sensor: [
// 		{x: 100, y: 100},
// 	]
// };

// draw(exdata);

function draw(data) {
	canvas.selectAll('.data-point').remove();

	for (let fern of data.ferns) {

		const x = (fern.location.x / 100) * 70;
		const y = (fern.location.y / 100) * 70;
		const fernx = width / 2 + x;
		const ferny = height / 2 + y;

		// console.log("Fern x y: " + x + " " + y);
		var chainAngle = 0;
		for (let chain of fern.leds) {
			var r = distBetweenLEDs;
			for (let led of chain) {
				// console.log("Ledr: " + led.r)
				const ledx = fernx + r*Math.cos(chainAngle);
				const ledy = ferny + r*Math.sin(chainAngle);
				canvas.append('circle')
					.attr('class', 'data-point')
					.attr('r', 3)
					.style('fill', 'rgb(' + led.R + ',' + led.G + ',' + led.B + ')')
					.attr('transform', 'translate(' + ledx + ', ' +	ledy + ')');

				r = r + distBetweenLEDs;
			}
			chainAngle += (2*Math.PI) / numChainsPerFern;
		}
	}

	for (let row of data.sensor) {
		const x = (row.x / 100) * 70;
		const y = (row.y / 100) * 70;
		// console.log(row);
		canvas.append('circle')
			.attr('class', 'data-point')
			.attr('r', 5)
			.style('fill', 'rgb(255, 0, 0)')
			.attr('transform', 'translate(' + ((width / 2) + x) + ', ' +
				((height / 2) + y) + ')');
	}
}


socket.addEventListener('message', function (event) {
	if (paused) {
		return;
	}

	var data = JSON.parse(event.data);
	draw(data);
});

canvas.on("click", function() {
	let coords = d3.mouse(this);
	socket.send(JSON.stringify({
		x: Math.floor((coords[0] - width/2) * 1.4285),
		y: Math.floor((coords[1] - height/2) * 1.4285)
	}));
});
