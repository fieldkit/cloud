import React, { PropTypes } from 'react'
import ViewportMercator from 'viewport-mercator-project'

const MapStateOverlay = (props) => {
  const { width, height, updateMap, currentDate, frameCount } = props
  console.log(frameCount)
  if (frameCount % 60 !== 59) return null

  // const mercator = ViewportMercator(props)
  // const { unproject } = mercator
  // var nw = unproject([0, 0])
  // var se = unproject([width, height])
  
  // updateMap(currentDate, [nw[0], nw[1], se[0], se[1]])
  return null
}

MapStateOverlay.propTypes = {
  width: PropTypes.number.isRequired,
  height: PropTypes.number.isRequired,
  updateMap: PropTypes.func.isRequired
}

export default MapStateOverlay