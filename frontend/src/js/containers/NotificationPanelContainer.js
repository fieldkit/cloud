
import { connect } from 'react-redux'
import NotificationPanel from '../components/NotificationPanel'
// import * as actions from '../actions'

const mapStateToProps = (state, ownProps) => {
  var expedition = state.expeditions[state.selectedExpedition]
  if (expedition) {
    return {
      posts: expedition.currentPosts,
      currentDate: expedition.currentDate,
      playback: expedition.playback
    }
  } else {
    return {}
  }
}

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
  }
}

const NotificationPanelContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(NotificationPanel)

export default NotificationPanelContainer
