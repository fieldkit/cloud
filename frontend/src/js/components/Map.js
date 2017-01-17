
import React, { PropTypes } from 'react'
import MapboxGL from 'react-map-gl'
import WebGLOverlay from './WebGLOverlay'
import DOMOverlay from './DOMOverlay'
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
      (playbackMode === 'forward' ? 500000 :
      playbackMode === 'fastForward' ? 5000000 :
      playbackMode === 'backward' ? -500000 :
      playbackMode === 'fastBackward' ? -5000000 : 
      0) / framesPerSecond
    const nextDate = Math.round(currentDate + dateDelta)
    updateDate(nextDate)
    requestAnimationFrame(this.tick)
  }

  componentDidMount () {
    requestAnimationFrame(this.tick)
  }

  shouldComponentUpdate (nextProps) {
    return this.props.currentDate !== nextProps.currentDate ||
      nextProps.playbackMode === 'pause'      
  }

  render () {
    const {
      viewport,
      setViewport,
      focusParticles,
      readingParticles,
      readingPath,
      focusedDocument
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
          <DOMOverlay
            focusedDocument={ focusedDocument }
          />
        </MapboxGL>
      </div>
    )
  }
}

Map.propTypes = {}

export default Map
