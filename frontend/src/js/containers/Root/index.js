
import { connect } from 'react-redux'
import * as actions from '../../actions'
import { createSelector } from 'reselect'

import Root from '../../components/Root'

const selectComputedData = createSelector(
  state => state.expeditions.get('currentExpedition'),
  state => {
    const currentExpedition = state.expeditions.get('currentExpedition')
    return state.expeditions.getIn(['expeditions', currentExpedition, 'expeditionFetching'])
  },
  state => {
    const currentExpedition = state.expeditions.get('currentExpedition')
    return state.expeditions.getIn(['expeditions', currentExpedition, 'documentsFetching'])
  },
  state => state.expeditions.get('documents'),
  (currentExpedition, expeditionFetching, documentsFetching, documents) => ({
    currentExpedition,
    expeditionFetching,
    documentsFetching,
    documents
  })
)

const mapStateToProps = (state, ownProps) => {
  return {
    ...selectComputedData(state)
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
