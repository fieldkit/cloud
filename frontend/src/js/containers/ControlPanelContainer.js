
import { connect } from 'react-redux'
import ControlPanel from '../components/ControlPanel'
import * as actions from '../actions'

const mapStateToProps = (state, ownProps) => {
  var props = {
    currentPage: state.currentPage,
    expeditionID: state.selectedExpedition,
    expeditions: state.expeditions
  }

  var expedition = state.expeditions[state.selectedExpedition]
  if (props.expeditionID) {
    props.currentDate = expedition.currentDate
    props.playback = expedition.playback
    props.mainFocus = expedition.mainFocus
    props.secondaryFocus = expedition.secondaryFocus
    props.zoom = expedition.zoom
    props.layout = expedition.layout
    props.viewport = {
      width: window.innerWidth,
      height: window.innerHeight,
      longitude: expedition.coordinates[0],
      latitude: expedition.coordinates[1],
      zoom: expedition.zoom
    }
  }

  return props
}

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    onYearChange: (value) => {
      return dispatch(actions.setExpedition(value))
    },
    onPlaybackChange: (value) => {
      return dispatch(actions.setControl('playback', value))
    },
    onMainFocusChange: (value) => {
      return dispatch(actions.setControl('mainFocus', value))
    },
    onSecondaryFocusChange: (value) => {
      return dispatch(actions.setControl('secondaryFocus', value))
    },
    onZoomChange: (value) => {
      return dispatch(actions.setControl('zoom', value))
    },
    onLayoutChange: (value) => {
      return dispatch(actions.setControl('layout', value))
    }
  }
}

const ControlPanelContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(ControlPanel)

export default ControlPanelContainer
