
import React, { PropTypes } from 'react'
import { dateToString } from '../../../utils.js'
import PlaybackSelector from './PlaybackSelector'
import FocusSelector from './FocusSelector'
import ZoomSelector from './ZoomSelector'

class ControlPanel extends React.Component {

    render() {
        const { activeExpedition, replay } = this.props
        const { selectPlaybackMode, selectFocusType, selectZoom} = this.props

        return (
            <div className="control-panel">
                <div className="control-panel_date-counter">
                    { dateToString(new Date(replay.now)) }
                </div>
                { replay.controlsVisible &&
                  <div>
                      <PlaybackSelector playbackMode={ replay.playbackMode } selectPlaybackMode={ selectPlaybackMode } />
                      <FocusSelector focus={ replay.focus } selectFocusType={ selectFocusType } />
                      <ZoomSelector zoom={ replay.zoom } selectZoom={ selectZoom } />
                  </div> }
            </div>
        )
    }
}

ControlPanel.propTypes = {
    activeExpedition: PropTypes.object.isRequired,
    replay: PropTypes.shape({
        controlsVisible: PropTypes.bool.isRequired,
        zoom: PropTypes.number.isRequired,
        now: PropTypes.number.isRequired,
        playbackMode: PropTypes.string.isRequired,
        focus: PropTypes.object.isRequired,
    }).isRequired,

    selectZoom: PropTypes.func.isRequired,
    selectPlaybackMode: PropTypes.func.isRequired,
    selectFocusType: PropTypes.func.isRequired,
}

export default ControlPanel
