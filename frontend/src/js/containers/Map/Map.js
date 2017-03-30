
import { connect } from 'react-redux'
import * as actions from '../../actions'
import { Vector3, Color } from 'three'
import ViewportMercator from 'viewport-mercator-project'
import { createSelector } from 'reselect'

import Map from '../../components/Map/Map'

const mapStateToProps = (state, ownProps) => {
  return {
    ...createSelector(
      state => state.expeditions.get('expeditions'),
      state => state.expeditions.get('viewport'),
      state => state.expeditions.get('documents'),
      state => state.expeditions.get('currentDocuments'),
      state => state.expeditions.get('currentDate'),
      state => state.expeditions.get('playbackMode'),
      state => state.expeditions.get('mousePosition'),
      state => state.expeditions.getIn(['focus', 'type']),
      (expeditions, _viewport, documents, currentDocumentIDs, currentDate, playbackMode, mousePosition, focusType) => {
        const viewport = _viewport.toJS()
        const currentDocuments = documents.filter(d => {
          return currentDocumentIDs.includes(d.get('id'))
        })
        let focusedDocument = null
        let focusDistance = 25
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
          viewport: _viewport,
          currentDocuments, 
          currentDate,
          playbackMode,
          focusParticles,
          readingParticles,
          readingPath,
          focusedDocument,
          focusType
        }
      }
    )(state)
  }
}

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    setViewport (viewport, manual) {
      return dispatch(actions.setViewport(viewport, manual))
    },
    updateDate (date) {
      return dispatch(actions.updateDate(date))
    },
    openLightbox (id) {
      return dispatch(actions.openLightbox(id))
    }
  }
}

const MapContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(Map)

export default MapContainer
