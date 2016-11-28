import React, {PropTypes} from 'react'
import * as d3 from 'd3'

const SightingGraph = ({width, height, sightings}) => {
  const margin = 0.42
  width = width - (margin * width)

  var x = d3.scaleLinear()
    .domain([0, d3.max(sightings)])
    .range([width, width * margin])

  var y = d3.scaleLinear()
    .domain([0, sightings.length - 1])
    .range([height, 0])

  return (
    <g>
      <path fill={'rgba(196,131,39,0.8)'} d={
        d3.area()
          .x0(width)
          .x1((d) => x(d))
          .y((d, i) => y(i))(sightings)
      }/>
    </g>
  )
}

SightingGraph.propTypes = {
  width: PropTypes.number.isRequired,
  height: PropTypes.number.isRequired,
  sightings: PropTypes.array
}

export default SightingGraph
