
import React, { PropTypes } from 'react'
import { dateToString } from '../../utils.js'

class ControlPanel extends React.Component {

  constructor (props) {
    super(props)
    this.state = {
    }
  }

  render () {

    const { currentDate, selectPlaybackMode } = this.props

    return (
      <div className="control-panel">
        <div className="control-panel_date-counter">
          { dateToString(new Date(currentDate)) }
        </div>
        <ul className="control-panel_playback-selector">
          <li
            className="control-panel_playback-selector_button"
            onClick={ () => selectPlaybackMode('fastBackward') }
          >
            fb
          </li>
          <li
            className="control-panel_playback-selector_button"
            onClick={ () => selectPlaybackMode('backward') }
          >
            b
          </li>
          <li
            className="control-panel_playback-selector_button"
            onClick={ () => selectPlaybackMode('pause') }
          >
            p
          </li>
          <li
            className="control-panel_playback-selector_button"
            onClick={ () => selectPlaybackMode('forward') }
          >
            f
          </li>
          <li
            className="control-panel_playback-selector_button"
            onClick={ () => selectPlaybackMode('fastForward') }
          >
            ff
          </li>
        </ul>
      </div>
    )
  }

}

ControlPanel.propTypes = {

}

export default ControlPanel