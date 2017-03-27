
import { connect } from 'react-redux'
import * as actions from '../../actions'
import I from 'immutable'
import { createSelector } from 'reselect'

import ExpeditionPanel from '../../components/Root/ExpeditionPanel'

const selectComputedData = createSelector(
  state => state.expeditions.get('currentExpedition'),
  state => state.expeditions.get('expeditionPanelOpen'),
  state => state.expeditions.get('expeditions'),
  state => state.expeditions.get('project'),
  (currentExpedition, expeditionPanelOpen, expeditions, project) => ({
    currentExpedition,
    expeditionPanelOpen,
    expeditions,
    project
  })
)

const mapStateToProps = (state, ownProps) => {
  return {
    ...selectComputedData(state)
  }
}

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    closeExpeditionPanel () {
      return dispatch(actions.closeExpeditionPanel())
    }
  }
}

const ExpeditionPanelContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(ExpeditionPanel)

export default ExpeditionPanelContainer
