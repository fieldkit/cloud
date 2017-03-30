
import React, { PropTypes } from 'react'
import MapboxGL from 'react-map-gl'
import WebGLOverlay from './WebGLOverlay'
import DOMOverlay from './DOMOverlay'
import { MAPBOX_ACCESS_TOKEN, MAPBOX_STYLE } from '../../constants/mapbox.js'
import { is, fromJS } from 'immutable'

class Map extends React.Component {

  constructor (props) {
    super(props)
    this.state = {}
    this.tick = this.tick.bind(this)
    this.onChangeViewport = this.onChangeViewport.bind(this)
  }

  tick (firstFrame) {
    const {
      currentDate,
      playbackMode,
      updateDate,
      focusType
    } = this.props
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

  onChangeViewport (viewport) {
    if (!is(this.props.viewport, fromJS(viewport))) {
      this.props.setViewport(viewport, true) 
    }
  }

  componentDidMount () {
    requestAnimationFrame(() => this.tick(true))
  }

  shouldComponentUpdate (props) {
    return !is(this.props.viewport, props.viewport) ||
      this.props.currentDate !== props.currentDate ||
      !is(this.props.focusedDocument, props.focusedDocument)
  }

  render () {
    const {
      setViewport,
      focusParticles,
      readingParticles,
      readingPath,
      focusedDocument,
      openLightbox
    } = this.props

    const viewport = this.props.viewport.toJS()

    return (
      <div id="map">
        <MapboxGL
          { ...viewport }
          mapStyle={ MAPBOX_STYLE }
          mapboxApiAccessToken={ MAPBOX_ACCESS_TOKEN }
          onChangeViewport={ this.onChangeViewport }
        >
          <WebGLOverlay
            { ...viewport }
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
