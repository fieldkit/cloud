
import React, { PropTypes } from 'react'
import ViewportMercator from 'viewport-mercator-project'
import I from 'immutable'
// import { lerp, parseDate } from '../utils'
import MapboxGL from 'react-map-gl'
import WebGLOverlay from './WebGLOverlay'
import DOMOverlay from './DOMOverlay'
import THREE from '../../vendors/react-three-renderer/node_modules/three'
import { MAPBOX_ACCESS_TOKEN, MAPBOX_STYLE } from '../constants/mapbox.js'

class Map extends React.Component {

  constructor (props) {
    super(props)
    this.tick = this.tick.bind(this)
  }

  tick () {
    const { currentDate, playbackMode, updateDate } = this.props
    const framesPerSecond = 60
    const dateDelta = 
      (playbackMode === 'forward' ? 100000 :
      playbackMode === 'fastForward' ? 1000000 :
      playbackMode === 'backward' ? -100000 :
      playbackMode === 'fastBackward' ? -1000000 : 
      0) / framesPerSecond
    const nextDate = Math.round(currentDate + dateDelta)
    updateDate(nextDate)
    requestAnimationFrame(this.tick)
  }

  componentDidMount () {
    requestAnimationFrame(this.tick)
  }

  render () {
    const {
      viewport,
      setViewport,
      focusParticles,
      readingParticles,
      readingPath
    } = this.props

    return (
      <div id="map">
        <MapboxGL
          { ...viewport }
          mapStyle={MAPBOX_STYLE}
          mapboxApiAccessToken={MAPBOX_ACCESS_TOKEN}
          onChangeViewport={viewport => {
            setViewport(viewport)
          }}
        >
          <WebGLOverlay
            { ...viewport }
            startDragLngLat={[0, 0]}
            redraw={this.redrawGLOverlay}
            focusParticles={ focusParticles }
            readingParticles={ readingParticles }
            readingPath={ readingPath }
          />
          {/*
            <DOMOverlay
              { ...viewport }
              startDragLngLat={[0, 0]}
              redraw={this.redrawDOMOverlay}
            />
          */}
        </MapboxGL>
      </div>
    )
  }

}

Map.propTypes = {

}

export default Map
