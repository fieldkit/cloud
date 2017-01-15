
import { connect } from 'react-redux'
import NotificationPanel from '../components/NotificationPanel'
import * as actions from '../actions'

const mapStateToProps = (state, ownProps) => {

  const currentDate = state.expeditions.get('currentDate')
  const currentDocuments = state.expeditions.get('documents')
    .filter(d => {
      return state.expeditions.get('currentDocuments').includes(d.get('id'))
    })

  return {
    currentDate,
    currentDocuments
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
