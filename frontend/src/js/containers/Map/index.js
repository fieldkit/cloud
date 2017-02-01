
import { connect } from 'react-redux'
import Map from '../../components/Map'
import * as actions from '../../actions'
import { Vector3, Color } from 'three'
import ViewportMercator from 'viewport-mercator-project'

const mapStateToProps = (state, ownProps) => {

  const expeditions = state.expeditions.get('expeditions')
  const viewport = state.expeditions.get('viewport').toJS()
  const currentDocuments = state.expeditions.get('documents')
    .filter(d => {
      return state.expeditions.get('currentDocuments').includes(d.get('id'))
    })
  const currentDate = state.expeditions.get('currentDate')
  const playbackMode = state.expeditions.get('playbackMode')
  const mousePosition = state.expeditions.get('mousePosition')
  let focusedDocument = null
  let focusDistance = 25

  // 
  // MAP BUFFER GEOMETRIES
  // 

  const particleCount = 1000
  const pathVertexCount = 200
  const { unproject } = ViewportMercator(viewport)
  const screenBounds = [[0, 0], [window.innerWidth, window.innerHeight]].map(unproject)
  
  // Focus particles
  const focusParticles = {
    position: new Float32Array(3),
    color: new Float32Array(4)
  }
  const radius = 15
  const x = window.innerWidth * ((viewport.longitude - screenBounds[0][0]) / (screenBounds[1][0] - screenBounds[0][0]))
  const y = window.innerHeight * ((viewport.latitude - screenBounds[0][1]) / (screenBounds[1][1] - screenBounds[0][1]))
  focusParticles.position[0] = x
  focusParticles.position[1] = y
  focusParticles.position[2] = radius
  focusParticles.color[0] = 242/255
  focusParticles.color[1] = 76/255
  focusParticles.color[2] = 45/255
  focusParticles.color[3] = 1

  // Readings particles and path
  const readingParticles = {
    position: new Float32Array(particleCount * 3),
    color: new Float32Array(particleCount * 4)
  }
  const readingPath = new Array(pathVertexCount)

  currentDocuments.toList().forEach((sighting, i) => {
    const position = sighting.getIn(['geometry', 'coordinates'])
    const radius = 15
    const x = window.innerWidth * ((sighting.getIn(['geometry', 'coordinates', 1]) - screenBounds[0][0]) / (screenBounds[1][0] - screenBounds[0][0]))
    const y = window.innerHeight * ((sighting.getIn(['geometry', 'coordinates', 0]) - screenBounds[0][1]) / (screenBounds[1][1] - screenBounds[0][1]))
    const color = new Color('#ffffff')
    
    readingParticles.position[i * 3 + 0] = x
    readingParticles.position[i * 3 + 1] = y
    readingParticles.position[i * 3 + 2] = radius
    readingParticles.color[i * 4 + 0] = color.r
    readingParticles.color[i * 4 + 1] = color.g
    readingParticles.color[i * 4 + 2] = color.b
    readingParticles.color[i * 4 + 3] = 1

    readingPath[i] = new Vector3(x, y, 0)

    const distanceToMouse = Math.sqrt(Math.pow(x - mousePosition.get(0), 2) + Math.pow(y - mousePosition.get(1), 2))
    if (distanceToMouse < focusDistance) {
      focusDistance = distanceToMouse
      focusedDocument = sighting
        .set('x', x)
        .set('y', y)
    }
  })

  // clear up unused particles
  readingParticles.position.fill(0, currentDocuments.size * 3, pathVertexCount * 3)
  readingParticles.color.fill(0, currentDocuments.size * 4, pathVertexCount * 4)
  readingPath.fill(readingPath[currentDocuments.size - 1], currentDocuments.size, pathVertexCount)

  return {
    expeditions,
    viewport,
    currentDocuments, 
    currentDate,
    playbackMode,
    focusParticles,
    readingParticles,
    readingPath,
    focusedDocument
  }
}

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    saveChangesAndResume () {
      return dispatch(actions.saveChangesAndResume())
    },
    cancelAction () {
      return dispatch(actions.cancelAction())
    },
    setViewport (viewport) {
      return dispatch(actions.setViewport(viewport))
    },
    updateDate (date) {
      return dispatch(actions.updateDate(date))
    }
  }
}

const MapContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(Map)

export default MapContainer
