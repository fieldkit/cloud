
import { connect } from 'react-redux'
import Root from '../../components/Root'
import * as actions from '../../actions'

const mapStateToProps = (state, ownProps) => {

  const currentExpedition = state.expeditions.get('currentExpedition')
  const expeditionFetching = state.expeditions.getIn(['expeditions', currentExpedition, 'expeditionFetching'])
  const documentsFetching = state.expeditions.getIn(['expeditions', currentExpedition, 'documentsFetching'])
  const documents = state.expeditions.get('documents')

  return {
    expeditionFetching,
    documentsFetching,
    currentExpedition,
    documents
  }
}

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    setMousePosition (x, y) {
      return dispatch(actions.setMousePosition(x, y))
    }
  }
}

const RootContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(Root)

export default RootContainer
