
import { connect } from 'react-redux'
import ControlPanel from '../../components/common/ControlPanel'
import * as actions from '../../actions'

const mapStateToProps = (state, ownProps) => {

  // const currentExpeditionID = state.expeditions.get('currentExpedition')
  // const expeditionName = state.expeditions.getIn(['expeditions', currentExpeditionID, 'name'])
  const currentDate = state.expeditions.get('currentDate')

  return {
    currentDate
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

const ControlPanelContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(ControlPanel)

export default ControlPanelContainer
