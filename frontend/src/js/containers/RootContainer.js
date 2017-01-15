
import { connect } from 'react-redux'
import Root from '../components/Root'
import * as actions from '../actions'

const mapStateToProps = (state, ownProps) => {

  const currentExpedition = state.expeditions.get('currentExpedition')
  const expeditionFetching = state.expeditions.getIn(['expeditions', currentExpedition, 'expeditionFetching'])
  const documentsFetching = state.expeditions.getIn(['expeditions', currentExpedition, 'documentsFetching'])

  return {
    expeditionFetching,
    documentsFetching
  }
}

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    saveChangesAndResume () {
      return dispatch(actions.saveChangesAndResume())
    },
    cancelAction () {
      return dispatch(actions.cancelAction())
    }
  }
}

const RootContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(Root)

export default RootContainer
