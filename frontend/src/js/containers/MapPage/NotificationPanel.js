
import { connect } from 'react-redux'
import * as actions from '../../actions'
import { createSelector } from 'reselect'

import NotificationPanel from '../../components/MapPage/NotificationPanel'

const mapStateToProps = (state, ownProps) => {
  return {
    ...createSelector(
      state => state.expeditions.get('currentDate'),
      state => state.expeditions.get('documents'),
      state => state.expeditions.get('currentDocuments'),
      (currentDate, documents, currentDocuments) => ({
        currentDate,
        currentDocuments: documents
          .filter(d => {
            return state.expeditions.get('currentDocuments').includes(d.get('id'))
          })
      })
    )(state)
  }
}

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    selectPlaybackMode (mode) {
      return dispatch(actions.selectPlaybackMode(mode))
    },
    cancelAction () {
      return dispatch(actions.cancelAction())
    }
  }
}

const NotificationPanelContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(NotificationPanel)

export default NotificationPanelContainer
