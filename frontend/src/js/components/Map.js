
import React, { PropTypes } from 'react'
import ViewportMercator from 'viewport-mercator-project'
import I from 'Immutable'
// import { lerp, parseDate } from '../utils'
import MapboxGL from 'react-map-gl'
import WebGLOverlay from './WebGLOverlay'
import THREE from '../react-three-renderer/node_modules/three'
import { MAPBOX_ACCESS_TOKEN, MAPBOX_STYLE } from '../constants/mapbox.js'

class Map extends React.Component {

  constructor (props) {
    super(props)
    this.redrawGLOverlay = this.redrawGLOverlay.bind(this)
    this.mapToScreen = this.mapToScreen.bind(this)
    this.renderSightings = this.renderSightings.bind(this)
  }

  redrawGLOverlay ({ unproject } ) {
    const screenBounds = [[0, 0], [window.innerWidth, window.innerHeight]].map(unproject)
    return (particles, paths) => {
      return {
        particles: {
          ...particles,
          sightings: this.renderSightings(particles.sightings, screenBounds, this.props.currentDocuments),
        }
      }
    }
  }

  mapToScreen (p, screenBounds) {
    return [
      window.innerWidth * ((p[0] - screenBounds[0][0]) / (screenBounds[1][0] - screenBounds[0][0])),
      window.innerHeight * ((p[1] - screenBounds[0][1]) / (screenBounds[1][1] - screenBounds[0][1]))
    ]
  }

  renderSightings (particleGeometry, screenBounds, sightings) {
    sightings.toList().forEach((sighting, i) => {
      const position = sighting.getIn(['geometry', 'coordinates'])
      const radius = 10
      const coords = this.mapToScreen([position.get(0), position.get(1)], screenBounds)
      const color = new THREE.Color('#ffffff')
      particleGeometry.position.array[i * 3 + 0] = coords[0]
      particleGeometry.position.array[i * 3 + 1] = coords[1]
      particleGeometry.position.array[i * 3 + 2] = radius * 2
      particleGeometry.color.array[i * 4 + 0] = color.r
      particleGeometry.color.array[i * 4 + 1] = color.g
      particleGeometry.color.array[i * 4 + 2] = color.b
      particleGeometry.color.array[i * 4 + 3] = 1
    })

    for (let i = sightings.size; i < particleGeometry.count; i++) {
      particleGeometry.position.array[i * 3 + 0] = 0
      particleGeometry.position.array[i * 3 + 1] = 0
      particleGeometry.position.array[i * 3 + 2] = 0
      particleGeometry.color.array[i * 4 + 0] = 0
      particleGeometry.color.array[i * 4 + 1] = 0
      particleGeometry.color.array[i * 4 + 2] = 0
      particleGeometry.color.array[i * 4 + 3] = 0
    }

    particleGeometry.position.needsUpdate = true
    particleGeometry.color.needsUpdate = true
    particleGeometry.data = sightings
    return particleGeometry
  }

  render () {
    const { viewport, setViewport, currentDocuments } = this.props

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
          />
        </MapboxGL>
      </div>
    )
  }

}

Map.propTypes = {

}

export default Map



