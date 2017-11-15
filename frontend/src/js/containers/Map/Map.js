
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

                // This is where all the geometry stuff for rendering documents on the map
                // TODO: This certainly could be optimized, results could be cached

                const viewport = _viewport.toJS()
                const currentDocuments = documents.filter(d => {
                    return currentDocumentIDs.includes(d.get('id'))
                }).sortBy(a => new Date(a.get("DateTime") || a.get("created_at")))
                const pathDocs = currentDocuments
                let focusedDocument = null
                let focusDistance = 25
                const particleCount = 1000
                const {unproject} = ViewportMercator(viewport)
                const screenBounds = [[0, 0], [window.innerWidth, window.innerHeight]].map(unproject)

                const documentTypes = []
                currentDocuments.forEach(d => {
                    if (documentTypes.indexOf(d.get('type')) === -1) documentTypes.push(d.get('type'))
                })

                // particles and path
                const particles = {}
                const readingPath = []
                const pointsPath = []
                documentTypes.forEach(id => {
                    particles[id] = {
                        position: new Float32Array(particleCount * 3),
                        color: new Float32Array(particleCount * 4)
                    }
                })

                pathDocs.toList().forEach((d, i) => {
                    const x = window.innerWidth * ((d.getIn(['geometry', 'coordinates', 1]) - screenBounds[0][0]) / (screenBounds[1][0] - screenBounds[0][0]))
                    const y = window.innerHeight * ((d.getIn(['geometry', 'coordinates', 0]) - screenBounds[0][1]) / (screenBounds[1][1] - screenBounds[0][1]))
                    const z = Math.abs(d.get("date") - currentDate)
                    pointsPath[i] = [x, y, z, d]
                })
                pathDocs.toList()
                    .filter(d => !d.get("user"))
                    .forEach((d, i) => {
                        const x = window.innerWidth * ((d.getIn(['geometry', 'coordinates', 1]) - screenBounds[0][0]) / (screenBounds[1][0] - screenBounds[0][0]))
                        const y = window.innerHeight * ((d.getIn(['geometry', 'coordinates', 0]) - screenBounds[0][1]) / (screenBounds[1][1] - screenBounds[0][1]))
                        const z = Math.abs(d.get("date") - currentDate)
                        readingPath[i] = [x, y, z, d]
                    })
                currentDocuments.toList()
                    .forEach((d, i) => {
                        const type = d.get('type')
                        const position = d.getIn(['geometry', 'coordinates'])
                        let radius = 15
                        const x = window.innerWidth * ((d.getIn(['geometry', 'coordinates', 1]) - screenBounds[0][0]) / (screenBounds[1][0] - screenBounds[0][0]))
                        const y = window.innerHeight * ((d.getIn(['geometry', 'coordinates', 0]) - screenBounds[0][1]) / (screenBounds[1][1] - screenBounds[0][1]))
                        let color,
                            s
                        let delta = Math.abs(d.get("date") - currentDate)
                        if (delta < 100000) {
                            radius = 15 + (202 * ((100000 - delta) / 100000))
                        }

                        if (d.get("user")) {
                            color = new Color('#00aced')
                        } else {
                            const speed = d.get("GPSSpeed") || 0
                            const r = Math.floor(speed)
                            color = new Color('#D0462C')
                            color.addScalar(speed / 2)
                        }

                        particles[type].position[i * 3 + 0] = x
                        particles[type].position[i * 3 + 1] = y
                        particles[type].position[i * 3 + 2] = radius
                        particles[type].color[i * 4 + 0] = color.r
                        particles[type].color[i * 4 + 1] = color.g
                        particles[type].color[i * 4 + 2] = color.b
                        particles[type].color[i * 4 + 3] = 1


                        const distanceToMouse = Math.sqrt(Math.pow(x - mousePosition.get(0), 2) + Math.pow(y - mousePosition.get(1), 2))
                        if (distanceToMouse < focusDistance) {
                            focusDistance = distanceToMouse
                            focusedDocument = d
                                .set('x', x)
                                .set('y', y)
                        }
                    })

                // clear up unused particles
                Object.keys(particles).forEach(type => {
                    const documentCount = currentDocuments
                        .filter(d => d.get('type') === type)
                        .size
                    particles[type].position.fill(0, documentCount * 3, particleCount * 3)
                    particles[type].color.fill(0, documentCount * 4, particleCount * 4)
                })

                return {
                    expeditions,
                    viewport: _viewport,
                    currentDocuments,
                    currentDate,
                    playbackMode,
                    particles,
                    focusedDocument,
                    focusType,
                    readingPath,
                    pointsPath
                }
            }
        )(state)
    }
}

const mapDispatchToProps = (dispatch, ownProps) => {
    return {
        setViewport(viewport, manual) {
            return dispatch(actions.setViewport(viewport, manual))
        },
        updateDate(date) {
            return dispatch(actions.updateDate(date))
        },
        openLightbox(id) {
            return dispatch(actions.openLightbox(id))
        }
    }
}

const MapContainer = connect(
    mapStateToProps,
    mapDispatchToProps
)(Map)

export default MapContainer
