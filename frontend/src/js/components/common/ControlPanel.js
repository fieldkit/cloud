
import React, { PropTypes } from 'react'
import { dateToString } from '../../utils.js'
import PlaybackSelector from './ControlPanel/PlaybackSelector'
import FocusSelector from './ControlPanel/FocusSelector'
import ZoomSelector from './ControlPanel/ZoomSelector'

class ControlPanel extends React.Component {

  render () {

    const { 
      currentDate,
      playbackMode,
      focus,
      zoom,
      selectPlaybackMode,
      selectFocusType,
      selectZoom
    } = this.props

    return (
      <div className="control-panel">
        <div className="control-panel_date-counter">
          { dateToString(new Date(currentDate)) }
        </div>
        <PlaybackSelector
          playbackMode={ playbackMode }
          selectPlaybackMode={ selectPlaybackMode }
        />
        <FocusSelector
          focus={ focus }
          selectFocusType={ selectFocusType }
        />
        <ZoomSelector
          zoom={ zoom }
          selectZoom={ selectZoom }
        />
      </div>
    )
  }

}

ControlPanel.propTypes = {

}

export default ControlPanel