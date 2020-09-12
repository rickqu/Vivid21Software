let width = 1000;
let height = 700;
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

const socket = new WebSocket('ws://127.0.0.1:9001/ws');

socket.addEventListener('message', function (event) {
	if (paused) {
		return;
	}

	var data = JSON.parse(event.data);
	canvas.selectAll('.data-point').remove();
	for (let row of data) {
		const rad = (row.a / 180) * Math.PI;
		const x = (row.x / 100) * 70;
		const y = (row.y / 100) * 70;
		console.log(row.s);
		canvas.append('circle')
			.attr('class', 'data-point')
			.attr('r', 3)
			.style('fill', row.color + (row.s/200) + ')')
			.attr('transform', 'translate(' + ((width / 2) + x) + ', ' +
				((height / 2) + y) + ')');
	}
});