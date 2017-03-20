
import { connect } from 'react-redux'
import ControlPanel from '../../components/common/ControlPanel'
import * as actions from '../../actions'

const mapStateToProps = (state, ownProps) => {

  // const currentExpeditionID = state.expeditions.get('currentExpedition')
  // const expeditionName = state.expeditions.getIn(['expeditions', currentExpeditionID, 'name'])
  const currentDate = state.expeditions.get('currentDate')
  const playbackMode = state.expeditions.get('playbackMode')
  const focus = state.expeditions.get('focus')

  return {
    currentDate,
    playbackMode,
    focus
  }
}

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    selectPlaybackMode (mode) {
      return dispatch(actions.selectPlaybackMode(mode))
    },
    selectFocusType (type) {
      return dispatch(actions.selectFocusType(type))
    },
    cancelAction () {
      return dispatch(actions.cancelAction())
    }
  }
}

const ControlPanelContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(ControlPanel)

export default ControlPanelContainer
