
import { connect } from 'react-redux'
import * as actions from '../../actions'
import { createSelector } from 'reselect'

import Header from '../../components/Root/Header'

const mapStateToProps = (state, ownProps) => {
  return {
    ...createSelector(
      state => state.expeditions.get('currentExpedition'),
      state => state.expeditions.getIn(['expeditions', state.expeditions.get('currentExpedition'), 'name']),
      state => state.expeditions.get('currentPage'),
      state => state.expeditions.getIn(["project","id"]),
      (currentExpeditionID, expeditionName, currentPage, projectID) => ({
        currentExpeditionID,
        expeditionName,
        currentPage,
        projectID
      })
    )(state)
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
    openExpeditionPanel () {
      return dispatch(actions.openExpeditionPanel())
    }
  }
}

const HeaderContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(Header)

export default HeaderContainer
