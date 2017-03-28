
import { connect } from 'react-redux'
import * as actions from '../../actions'
import { createSelector } from 'reselect'

import Root from '../../components/Root/Root'

const mapStateToProps = (state, ownProps) => {
  return {
    ...createSelector(
      state => state.expeditions.get('currentExpedition'),
      state => state.expeditions.getIn(['expeditions', state.expeditions.get('currentExpedition'), 'expeditionFetching']),
      state => state.expeditions.getIn(['expeditions', state.expeditions.get('currentExpedition'), 'documentsFetching']),
      state => state.expeditions.get('documents'),
      (currentExpedition, expeditionFetching, documentsFetching, documents) => ({
        currentExpedition,
        expeditionFetching,
        documentsFetching,
        documents
      })
    )(state)
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
