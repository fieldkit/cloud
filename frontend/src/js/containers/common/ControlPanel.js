
import { connect } from 'react-redux'
import * as actions from '../../actions'
import { createSelector } from 'reselect'

import ControlPanel from '../../components/common/ControlPanel/ControlPanel'

const mapStateToProps = (state, ownProps) => {
    return {
        ...createSelector(
            state => state.expeditions.get('currentExpedition'),
            state => state.expeditions.get('currentDate'),
            state => state.expeditions.get('playbackMode'),
            state => state.expeditions.get('focus'),
            state => state.expeditions.getIn(['viewport', 'zoom']),
            state => state.expeditions.get('showSensors'),
            (currentExpeditionID, currentDate, playbackMode, focus, zoom, showSensors) => ({
                activeExpedition: {
                    id: currentExpeditionID,
                },
                replay: {
                    controlsVisible: true,
                    now: new Date(currentDate).getTime(),
                    playbackMode,
                    focus,
                    zoom,
                },
                currentExpeditionID,
                currentDate,
                playbackMode,
                focus,
                zoom,
            })
        )(state)
    }
}

const mapDispatchToProps = (dispatch, ownProps) => {
    return {
        selectPlaybackMode(mode) {
            return dispatch(actions.selectPlaybackMode(mode))
        },
        selectFocusType(type) {
            return dispatch(actions.selectFocusType(type))
        },
        selectZoom(zoom) {
            return dispatch(actions.selectZoom(zoom))
        },
    }
}

const ControlPanelContainer = connect(
    mapStateToProps,
    mapDispatchToProps
)(ControlPanel)

export default ControlPanelContainer
