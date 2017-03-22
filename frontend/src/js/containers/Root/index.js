
import { connect } from 'react-redux'
import Root from '../../components/Root'
import * as actions from '../../actions'
import I from 'immutable'

const mapStateToProps = (state, ownProps) => {

  const currentExpedition = state.expeditions.get('currentExpedition')
  const expeditionFetching = state.expeditions.getIn(['expeditions', currentExpedition, 'expeditionFetching'])
  const documentsFetching = state.expeditions.getIn(['expeditions', currentExpedition, 'documentsFetching'])
  const documents = state.expeditions.get('documents')
  const expeditionPanelOpen = state.expeditions.get('expeditionPanelOpen')

  const expeditions = state.expeditions.get('expeditions')
  const project = I.fromJS({id: 'eric', name: 'eric'})

  return {
    expeditionFetching,
    documentsFetching,
    currentExpedition,
    expeditionPanelOpen,
    documents,
    expeditions,
    project
  }
}

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    setMousePosition (x, y) {
      return dispatch(actions.setMousePosition(x, y))
    },
    closeExpeditionPanel () {
      return dispatch(actions.closeExpeditionPanel())
    }
  }
}

const RootContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(Root)

export default RootContainer
