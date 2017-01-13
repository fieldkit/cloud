
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
    // this.renderAmbitGeo = this.renderAmbitGeo.bind(this)
    this.renderSightings = this.renderSightings.bind(this)
    this.renderMembers = this.renderMembers.bind(this)
  }

  redrawGLOverlay ({ unproject } ) {
    const screenBounds = [[0, 0], [window.innerWidth, window.innerHeight]].map(unproject)
    return (particles, paths) => {
      const { expedition, viewport } = this.props
      const { geoBounds } = viewport
      const west = geoBounds[0] + (geoBounds[0] - geoBounds[2]) * 0.25
      const north = geoBounds[1] + (geoBounds[1] - geoBounds[3]) * 0.25
      const east = geoBounds[2] + (geoBounds[2] - geoBounds[0]) * 0.25
      const south = geoBounds[3] + (geoBounds[3] - geoBounds[1]) * 0.25
      const gb = [west, north, east, south]

      // console.log('aga', particles.sightings)

      const sightings = this.props.currentDocuments

      return {
        particles: {
          ...particles,
          sightings: this.renderSightings(particles.sightings, screenBounds, sightings, gb),
          // members: this.renderMembers(particles.members, screenBounds, expedition, gb)
        },
        // paths: {
        //   ambitGeo: this.renderAmbitGeo(paths.ambitGeo, screenBounds, expedition, gb)
        // },
      }
    }
  }

  mapToScreen (p, screenBounds) {
    return [
      0 + (window.innerWidth - 0) * ((p[0] - screenBounds[0][0]) / (screenBounds[1][0] - screenBounds[0][0])),
      0 + (window.innerHeight - 0) * ((p[1] - screenBounds[0][1]) / (screenBounds[1][1] - screenBounds[0][1]))
    ]
  }

  // renderAmbitGeo (pathGeometry, screenBounds, expedition, gb) {
  //   const checkGeoBounds = (p, gb) => {
  //     return p[0] >= gb[0] && p[0] < gb[2] && p[1] >= gb[3] && p[1] < gb[1]
  //   }

  //   return expedition.currentAmbits.map(route => {

  //     const vertices = route.coordinates
  //       .filter((p, i) => {
  //         if (route.coordinates[i - 1] && checkGeoBounds(route.coordinates[i - 1], gb)) return true
  //         if (checkGeoBounds(route.coordinates[i], gb)) return true
  //         if (route.coordinates[i + 1] && checkGeoBounds(route.coordinates[i + 1], gb)) return true
  //         return false
  //       })
  //       .map(p => {
  //         return this.mapToScreen(p, screenBounds)
  //       })
  //       .map((p, i) => {
  //         return new THREE.Vector3(p[0], p[1], 0)
  //       })

  //     var lastVertex = vertices[vertices.length - 1].clone()
  //     for (let i = vertices.length; i < 200; i ++) {
  //       vertices[i] = lastVertex
  //     }

  //     return {
  //       color: route.color,
  //       vertices
  //     }
  //   })
  // }

  renderSightings (particleGeometry, screenBounds, sightings, gb) {
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

  renderMembers (geometry, screenBounds, expedition, gb) {
    if (!this.state.members) return geometry
    const members = Object.keys(this.state.members).map(name => {
      const member = this.state.members[name]
      const position = this.mapToScreen(member.coordinates, screenBounds)
      return {
        name,
        position
      }
    })
    return members
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



