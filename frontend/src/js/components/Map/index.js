
import React, { PropTypes } from 'react'
import MapboxGL from 'react-map-gl'
import WebGLOverlay from './WebGLOverlay'
import DOMOverlay from './DOMOverlay'
import { MAPBOX_ACCESS_TOKEN, MAPBOX_STYLE } from '../../constants/mapbox.js'

class Map extends React.Component {

  constructor (props) {
    super(props)
    this.state = {
      viewport: {
        longitude: 0,
        latitude: 0,
        zoom: 15
      }
    }
    this.tick = this.tick.bind(this)
  }

  tick (firstFrame) {
    const { currentDate, playbackMode, updateDate } = this.props
    const framesPerSecond = 60
    const dateDelta = 
      (playbackMode === 'forward' ? 500000 :
      playbackMode === 'fastForward' ? 5000000 :
      playbackMode === 'backward' ? -500000 :
      playbackMode === 'fastBackward' ? -5000000 : 
      0) / framesPerSecond
    const nextDate = Math.round(currentDate + dateDelta)
    if (firstFrame || dateDelta !== 0) updateDate(nextDate)
    requestAnimationFrame(() => this.tick(false))
  }

  componentDidMount () {
    requestAnimationFrame(() => this.tick(true))
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
      focusedDocument,
      openLightbox
    } = this.props

    return (
      <div id="map">
        <MapboxGL
          { ...viewport }
          mapStyle={MAPBOX_STYLE}
          mapboxApiAccessToken={MAPBOX_ACCESS_TOKEN}
          onChangeViewport={ viewport => setViewport(viewport, true) }
        >
          <WebGLOverlay
            { ...viewport }
            startDragLngLat={ [0, 0] }
            redraw={ this.redrawGLOverlay }
            focusParticles={ focusParticles }
            readingParticles={ readingParticles }
            readingPath={ readingPath }
          />
          <DOMOverlay
            focusedDocument={ focusedDocument }
            openLightbox={ openLightbox }
          />
        </MapboxGL>
      </div>
    )
  }
}

Map.propTypes = {}

export default Map
