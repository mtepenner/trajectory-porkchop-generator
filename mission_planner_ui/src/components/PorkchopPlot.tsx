import React, { useEffect, useRef } from 'react';
import * as d3 from 'd3';
import { contours } from 'd3-contour';

interface Props {
  grid: number[][];
  depAxis: number[];
  arrAxis: number[];
  onCellClick: (dep: number, arr: number) => void;
}

const PorkchopPlot: React.FC<Props> = ({ grid, depAxis, arrAxis, onCellClick }) => {
  const svgRef = useRef<SVGSVGElement>(null);

  useEffect(() => {
    if (!svgRef.current || !grid.length) return;

    const width = 480;
    const height = 400;
    const margin = { top: 20, right: 20, bottom: 50, left: 60 };

    const svg = d3.select(svgRef.current)
      .attr('width', width)
      .attr('height', height);

    svg.selectAll('*').remove();

    const flat = grid.flat().filter((v) => isFinite(v));
    const [dMin, dMax] = d3.extent(flat) as [number, number];

    const xScale = d3.scaleLinear()
      .domain([depAxis[0], depAxis[depAxis.length - 1]])
      .range([margin.left, width - margin.right]);

    const yScale = d3.scaleLinear()
      .domain([arrAxis[0], arrAxis[arrAxis.length - 1]])
      .range([height - margin.bottom, margin.top]);

    const colorScale = d3.scaleSequential(d3.interpolateYlOrRd).domain([dMin, dMax]);

    const cellW = (width - margin.left - margin.right) / grid.length;
    const cellH = (height - margin.top - margin.bottom) / grid[0].length;

    const g = svg.append('g');

    grid.forEach((row, i) => {
      row.forEach((val, j) => {
        if (!isFinite(val)) return;
        g.append('rect')
          .attr('x', xScale(depAxis[i]))
          .attr('y', yScale(arrAxis[j]) - cellH)
          .attr('width', cellW)
          .attr('height', cellH)
          .attr('fill', colorScale(val))
          .style('cursor', 'pointer')
          .on('click', () => onCellClick(depAxis[i], arrAxis[j]));
      });
    });

    svg.append('g').attr('transform', `translate(0,${height - margin.bottom})`)
      .call(d3.axisBottom(xScale).ticks(6))
      .append('text')
      .attr('x', (width - margin.left - margin.right) / 2 + margin.left)
      .attr('y', 40)
      .attr('fill', '#ccc')
      .text('Departure MJD');

    svg.append('g').attr('transform', `translate(${margin.left},0)`)
      .call(d3.axisLeft(yScale).ticks(6))
      .append('text')
      .attr('transform', 'rotate(-90)')
      .attr('x', -(height / 2))
      .attr('y', -45)
      .attr('fill', '#ccc')
      .text('Arrival MJD');

  }, [grid, depAxis, arrAxis, onCellClick]);

  return (
    <div style={{ background: '#111', borderRadius: 8, padding: '0.5rem' }}>
      <h3 style={{ color: '#7af', margin: '0 0 0.5rem 0' }}>ΔV Porkchop Plot (km/s)</h3>
      <svg ref={svgRef} />
    </div>
  );
};

export default PorkchopPlot;
