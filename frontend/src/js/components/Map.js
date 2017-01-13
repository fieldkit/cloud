
import React, { PropTypes } from 'react'
import ViewportMercator from 'viewport-mercator-project'
// import { lerp, parseDate } from '../utils'
import MapboxGL from 'react-map-gl'
// import WebGLOverlay from './WebGLOverlay'
// import THREE from '../react-three-renderer/node_modules/three'

class Map extends React.Component {

  constructor (props) {
    super(props)
    this.state = {
      animate: false,
      viewport: {
        latitude: -18.5699229,
        longitude: 22.115456,
        zoom: 4,
        width: window.innerWidth,
        height: window.innerHeight,
        startDragLngLat: null,
        isDragging: false
      }
    }
  }

  render () {
    const { viewport } = this.state
    const MAPBOX_ACCESS_TOKEN = 'pk.eyJ1IjoiaWFhYWFuIiwiYSI6ImNpbXF1ZW4xOTAwbnl3Ymx1Y2J6Mm5xOHYifQ.6wlNzSdcTlonLBH-xcmUdQ'
    const MAPBOX_STYLE = 'mapbox://styles/mapbox/satellite-v9?format=jpg70'

    console.log(viewport)

    return (
      <div id="map">
        <MapboxGL
          latitude={-18.5699229}
          longitude={22.115456}
          zoom={4}
          width={window.innerWidth}
          height={window.innerHeight}
          mapStyle={MAPBOX_STYLE}
          mapboxApiAccessToken={MAPBOX_ACCESS_TOKEN}
        >
        </MapboxGL>
      </div>
    )
  }

}

Map.propTypes = {

}

export default Map



  // @autobind
  // tick (pastFrameDate) {
  //   let newState = {}
  //   const speedFactor = (Date.now() - pastFrameDate) / (1000 / 60)
  //   const currentFrameDate = Date.now()
  //   const {expeditionID, animate, expedition, fetchDay, setControl, isFetching, updateMap, initialPage} = this.props
  //   var b1, b2
  //   if (animate && !isFetching && location.pathname === '/map' || location.pathname === '/') {
  //     // increment time
  //     var dateOffset = 0
  //     var forward = expedition.playback === 'fastForward' || expedition.playback === 'forward' || expedition.playback === 'pause'
  //     if (this.state.beaconIndex === (forward ? 0 : 1) || this.state.beaconIndex === (forward ? d3.values(this.state.day.beacons).length - 2 : d3.values(this.state.day.beacons).length - 1)) {
  //       var offset = this.state.timeToNextBeacon > 0 ? Math.min(100000, this.state.timeToNextBeacon + 1) : 100000
  //       if (expedition.playback === 'fastBackward' || expedition.playback === 'backward') dateOffset = -1 * offset
  //       if (expedition.playback === 'forward' || expedition.playback === 'fastForward') dateOffset = offset
  //     } else {
  //       if (expedition.playback === 'fastBackward') dateOffset = -25000
  //       if (expedition.playback === 'backward') dateOffset = -4000
  //       if (expedition.playback === 'forward') dateOffset = 4000
  //       if (expedition.playback === 'fastForward') dateOffset = 25000
  //     }
  //     var currentDate = new Date(Math.min(expedition.end.getTime() - 1, (Math.max(expedition.start.getTime() + 1, this.state.currentDate.getTime() + dateOffset))))

  //     // pause playback if time reaches beginning or end
  //     if ((currentDate.getTime() === expedition.end.getTime() - 1 && (expedition.playback === 'forward' || expedition.playback === 'fastForward')) || (currentDate.getTime() === expedition.start.getTime() + 1 && (expedition.playback === 'backward' || expedition.playback === 'fastBackward'))) setControl('playback', 'pause')

  //     // checks current day
  //     var currentDay = Math.floor((currentDate.getTime() - expedition.start.getTime()) / (1000 * 3600 * 24))
  //     if (currentDay !== this.state.currentDay) {
  //       // new day
  //       fetchDay(currentDate)
  //     }

  //     // look for most current beacon
  //     const day = expedition.days[currentDay]
  //     var beacons = d3.values(day.beacons).sort((a, b) => {
  //       return parseDate(a.properties.DateTime).getTime() - parseDate(b.properties.DateTime).getTime()
  //     })
  //     var beaconCount = beacons.length
  //     var beaconIndex
  //     var timeToNextBeacon = 0
  //     var ratioBetweenBeacons = 0
  //     if (expedition.playback === 'forward' || expedition.playback === 'fastForward' || expedition.playback === 'pause') {
  //       for (var i = 0; i < beaconCount - 1; i++) {
  //         b1 = parseDate(beacons[i].properties.DateTime).getTime()
  //         b2 = parseDate(beacons[i + 1].properties.DateTime).getTime()
  //         if (currentDate.getTime() >= b1 && currentDate.getTime() < b2) {
  //           beaconIndex = i
  //           timeToNextBeacon = b2 - currentDate.getTime()
  //           ratioBetweenBeacons = (currentDate.getTime() - b1) / (b2 - b1)
  //           break
  //         }
  //       }
  //       if (beaconIndex < 0) beaconIndex = beaconCount - 1
  //     } else {
  //       for (i = beaconCount - 1; i > 0; i--) {
  //         b1 = parseDate(beacons[i].properties.DateTime).getTime()
  //         b2 = parseDate(beacons[i - 1].properties.DateTime).getTime()
  //         if (currentDate.getTime() <= b1 && currentDate.getTime() > b2) {
  //           beaconIndex = i
  //           timeToNextBeacon = currentDate.getTime() - b2
  //           ratioBetweenBeacons = (currentDate.getTime() - b1) / (b2 - b1)
  //           break
  //         }
  //       }
  //       if (beaconIndex < 0) beaconIndex = 0
  //     }
  //     // set map coordinates to current beacon
  //     var currentBeacon = beacons[beaconIndex + (forward ? 0 : 0)]
  //     var nextBeacon = beacons[beaconIndex + (forward ? 1 : -1)]
  //     var coordinates = [
  //       lerp(currentBeacon.geometry.coordinates[0], nextBeacon.geometry.coordinates[0], ratioBetweenBeacons),
  //       lerp(currentBeacon.geometry.coordinates[1], nextBeacon.geometry.coordinates[1], ratioBetweenBeacons)
  //     ]

  //      // look for most current ambit_geo
  //     const members = { ...expedition.members }
  //     Object.keys(members).forEach(memberID => {
  //       var member = members[memberID]
  //       var ambits = d3.values(expedition.featuresByMember[memberID][currentDay]).sort((a, b) => {
  //         return parseDate(a.properties.DateTime).getTime() - parseDate(b.properties.DateTime).getTime()
  //       })
  //       var ambitCount = ambits.length
  //       var ambitIndex = -1
  //       var ratioBetweenAmbits = 0
  //       if (expedition.playback === 'forward' || expedition.playback === 'fastForward' || expedition.playback === 'pause') {
  //         for (var i = 0; i < ambitCount - 1; i++) {
  //           b1 = parseDate(ambits[i].properties.DateTime).getTime()
  //           b2 = parseDate(ambits[i + 1].properties.DateTime).getTime()
  //           if (currentDate.getTime() >= b1 && currentDate.getTime() < b2) {
  //             ambitIndex = i
  //             ratioBetweenAmbits = (currentDate.getTime() - b1) / (b2 - b1)
  //             break
  //           }
  //         }
  //         if (ambitIndex < 0) {
  //           ambitIndex = ambitCount - 2
  //           ratioBetweenAmbits = 1
  //         }
  //       } else {
  //         for (i = ambitCount - 1; i > 0; i--) {
  //           b1 = parseDate(ambits[i].properties.DateTime).getTime()
  //           b2 = parseDate(ambits[i - 1].properties.DateTime).getTime()
  //           if (currentDate.getTime() <= b1 && currentDate.getTime() > b2) {
  //             ambitIndex = i
  //             ratioBetweenAmbits = (currentDate.getTime() - b1) / (b2 - b1)
  //             break
  //           }
  //         }
  //         if (ambitIndex < 0) {
  //           ambitIndex = 1
  //           ratioBetweenAmbits = 1
  //         }
  //       }
  //       // set member coordinates
  //       var currentID = ambitIndex
  //       var nextID = ambitIndex + (forward ? 1 : -1)
  //       if (currentID >= 0 && currentID < ambits.length && nextID >= 0 && nextID < ambits.length) {
  //         var currentAmbits = ambits[currentID]
  //         var nextAmbit = ambits[nextID]
  //         member.coordinates = [
  //           lerp(currentAmbits.geometry.coordinates[0], nextAmbit.geometry.coordinates[0], ratioBetweenAmbits),
  //           lerp(currentAmbits.geometry.coordinates[1], nextAmbit.geometry.coordinates[1], ratioBetweenAmbits)
  //         ]
  //       } else {
  //         member.coordinates = [-180, 90]
  //       }
  //     })

  //     var zoom = lerp(this.state.viewport.zoom, this.state.viewport.targetZoom, Math.pow(this.state.viewport.zoom / this.state.viewport.targetZoom, 2) / 250 * speedFactor)
  //     // if (!(initialPage === '/' || initialPage === '/map') || (!this.state.contentActive && this.props.contentActive)) zoom = this.state.viewport.targetZoom
  //     if (!(initialPage === '/' || initialPage === '/map')) zoom = this.state.viewport.targetZoom

  //     newState = {
  //       ...newState,
  //       currentDate,
  //       animate,
  //       currentDay,
  //       day,
  //       beaconIndex,
  //       timeToNextBeacon,
  //       members,
  //       contentActive: this.props.contentActive,
  //       viewport: {
  //         ...this.state.viewport,
  //         longitude: coordinates[0],
  //         latitude: coordinates[1],
  //         zoom: zoom
  //       }
  //     }

  //     if (this.state.frameCount % 60 === 0) {
  //       const { unproject } = ViewportMercator({ ...this.state.viewport })
  //       const nw = unproject([0, 0])
  //       const se = unproject([window.innerWidth, window.innerHeight])
  //       const viewGeoBounds = [nw[0], nw[1], se[0], se[1]]
  //       updateMap(this.state.currentDate, [this.state.viewport.longitude, this.state.viewport.latitude], viewGeoBounds, this.state.viewport.zoom, expeditionID)
  //     }
  //   }

  //   this.setState({
  //     ...this.state,
  //     ...newState,
  //     animate,
  //     frameCount: this.state.frameCount + 1
  //   })

  //   requestAnimationFrame(() => { this.tick(currentFrameDate) })
  // }


  // componentWillReceiveProps (nextProps) {
  //   const {animate, expedition, mapStateNeedsUpdate} = nextProps
  //   // console.log('new', animate, this.state.animate)
  //   if (animate) {
  //     const currentDate = expedition.currentDate
  //     // note: currentDay has a 1 day offset with API expeditionDay, which starts at 1
  //     const currentDay = Math.floor((currentDate.getTime() - expedition.start.getTime()) / (1000 * 3600 * 24))
  //     const day = expedition.days[currentDay]

  //     if (mapStateNeedsUpdate) {
  //       this.state.currentDate = currentDate
  //       this.state.currentDay = currentDay
  //       this.state.day = day
  //       this.state.frameCount = 0
  //     }

  //     if (!this.state.animate) {
  //       this.state.animate = animate
  //       this.state.viewport = {
  //         ...this.state.viewport,
  //         zoom: expedition.initialZoom,
  //         targetZoom: expedition.targetZoom
  //       }
  //       // console.log('starting animation')
  //       this.tick(Math.round(Date.now() - (1000 / 60)))
  //     }
  //   }
  // }

  // @autobind
  // redrawGLOverlay ({ unproject } ) {
  //   const screenBounds = [[0, 0], [window.innerWidth, window.innerHeight]].map(unproject)
  //   return (particles, paths) => {
  //     const { expedition } = this.props
  //     const { currentGeoBounds } = expedition
  //     const west = currentGeoBounds[0] + (currentGeoBounds[0] - currentGeoBounds[2]) * 0.25
  //     const north = currentGeoBounds[1] + (currentGeoBounds[1] - currentGeoBounds[3]) * 0.25
  //     const east = currentGeoBounds[2] + (currentGeoBounds[2] - currentGeoBounds[0]) * 0.25
  //     const south = currentGeoBounds[3] + (currentGeoBounds[3] - currentGeoBounds[1]) * 0.25
  //     const gb = [west, north, east, south]

  //     if (expedition.zoom < 14) {
  //       return {
  //         particles,
  //         paths
  //       }
  //     } else {
  //       return {
  //         particles: {
  //           ...particles,
  //           // pictures360: this.render360Images(particles.pictures360, screenBounds, expedition, gb),
  //           sightings: this.renderSightings(particles.sightings, screenBounds, expedition, gb),
  //           members: this.renderMembers(particles.members, screenBounds, expedition, gb)
  //         },
  //         paths: {
  //           ambitGeo: this.renderAmbitGeo(paths.ambitGeo, screenBounds, expedition, gb)
  //         },
  //       }
  //     }
  //   }
  // }

  // @autobind
  // mapToScreen (p, screenBounds) {
  //   return [
  //     0 + (window.innerWidth - 0) * ((p[0] - screenBounds[0][0]) / (screenBounds[1][0] - screenBounds[0][0])),
  //     0 + (window.innerHeight - 0) * ((p[1] - screenBounds[0][1]) / (screenBounds[1][1] - screenBounds[0][1]))
  //   ]
  // }

  // @autobind
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

  // @autobind
  // render360Images (particleGeometry, screenBounds, expedition, gb) {

  //   const images = expedition.current360Images
  //     .filter(image => {
  //       const coords = image.geometry.coordinates
  //       return coords[0] >= gb[0] && coords[0] < gb[2] && coords[1] >= gb[3] && coords[1] < gb[1]
  //     })

  //   for (var i = 0; i < particleGeometry.count; i++) {
  //     const image = images[i]
  //     if (image) {
  //       const coords = this.mapToScreen([image.geometry.coordinates[0], image.geometry.coordinates[1]], screenBounds)
  //       particleGeometry.position.array[i * 3 + 0] = coords[0]
  //       particleGeometry.position.array[i * 3 + 1] = coords[1]
  //       particleGeometry.position.array[i * 3 + 2] = 0
  //       particleGeometry.color.array[i * 4 + 0] = 0.62
  //       particleGeometry.color.array[i * 4 + 1] = 0.6
  //       particleGeometry.color.array[i * 4 + 2] = 0.7
  //       particleGeometry.color.array[i * 4 + 3] = 1
  //     } else {
  //       particleGeometry.position.array[i * 3 + 0] = 0
  //       particleGeometry.position.array[i * 3 + 1] = 0
  //       particleGeometry.position.array[i * 3 + 2] = 0
  //       particleGeometry.color.array[i * 4 + 0] = 0
  //       particleGeometry.color.array[i * 4 + 1] = 0
  //       particleGeometry.color.array[i * 4 + 2] = 0
  //       particleGeometry.color.array[i * 4 + 3] = 0
  //     }
  //   }

  //   particleGeometry.position.needsUpdate = true
  //   particleGeometry.color.needsUpdate = true
  //   particleGeometry.data = images
  //   return particleGeometry
  // }

  // @autobind
  // renderSightings (particleGeometry, screenBounds, expedition, gb) {
  //   const sightings = expedition.currentSightings
  //     .filter((sighting, i) => {
  //       const { position } = sighting
  //       return position.x >= gb[0] && position.x < gb[2] && position.y >= gb[3] && position.y < gb[1]
  //     })

  //   for (var i = 0; i < particleGeometry.count; i++) {
  //     const sighting = sightings[i]
  //     if (sighting) {
  //       const { position, radius } = sighting
  //       const coords = this.mapToScreen([position.x, position.y], screenBounds)
  //       const color = new THREE.Color(sighting.color)
  //       particleGeometry.position.array[i * 3 + 0] = coords[0]
  //       particleGeometry.position.array[i * 3 + 1] = coords[1]
  //       particleGeometry.position.array[i * 3 + 2] = radius * 2
  //       particleGeometry.color.array[i * 4 + 0] = color.r
  //       particleGeometry.color.array[i * 4 + 1] = color.g
  //       particleGeometry.color.array[i * 4 + 2] = color.b
  //       particleGeometry.color.array[i * 4 + 3] = 1
  //     } else {
  //       particleGeometry.position.array[i * 3 + 0] = 0
  //       particleGeometry.position.array[i * 3 + 1] = 0
  //       particleGeometry.position.array[i * 3 + 2] = 0
  //       particleGeometry.color.array[i * 4 + 0] = 0
  //       particleGeometry.color.array[i * 4 + 1] = 0
  //       particleGeometry.color.array[i * 4 + 2] = 0
  //       particleGeometry.color.array[i * 4 + 3] = 0
  //     }
  //   }

  //   particleGeometry.position.needsUpdate = true
  //   particleGeometry.color.needsUpdate = true
  //   particleGeometry.data = sightings
  //   return particleGeometry
  // }

  // @autobind
  // renderMembers (geometry, screenBounds, expedition, gb) {
  //   if (!this.state.members) return geometry
  //   const members = Object.keys(this.state.members).map(name => {
  //     const member = this.state.members[name]
  //     const position = this.mapToScreen(member.coordinates, screenBounds)
  //     return {
  //       name,
  //       position
  //     }
  //   })
  //   return members
  // }


