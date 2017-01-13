
import { connect } from 'react-redux'
import Map from '../components/Map'
import * as actions from '../actions'

const mapStateToProps = (state, ownProps) => {

  const expeditions = state.expeditions.get('expeditions')
  const viewport = state.expeditions.get('viewport').toJS()
  const currentDocuments = state.expeditions.get('documents')
    .filter(d => {
      return state.expeditions.get('currentDocuments').includes(d.get('id'))
    })

  return {
    expeditions,
    viewport,
    currentDocuments
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
    }
  }
}

const MapContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(Map)

export default MapContainer
