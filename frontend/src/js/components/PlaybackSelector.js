
import React, {PropTypes} from 'react'
import autobind from 'autobind-decorator'

class PlaybackSelector extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
      mouseOver: ''
    }
  }

  @autobind
  onMouseMove (e) {
    this.setState({
      ...this.state,
      mouseOver: e.nativeEvent.target.className
    })
  }

  @autobind
  onMouseOut (e) {
    this.setState({
      ...this.state,
      mouseOver: ''
    })
  }

  render () {
    const { mode, onPlaybackChange } = this.props
    const { mouseOver } = this.state
    const onMouseMove = this.onMouseMove
    var types = ['fastBackward', 'backward', 'pause', 'forward', 'fastForward']
    var buttons = types.map(function (s, i) {
      var className = 'playbackButton ' + (s === mode ? 'active' : 'inactive')

      var src = () => {
        var url = '/static/img/icon-' + s
        if (s === mode || mouseOver === s) url += '-hover'
        return url + '.png'
      }

      return (
        <li className={className} key={i} onClick={() => onPlaybackChange(s)} onMouseMove={onMouseMove}>
          <img className={s} width="16" height="16" src={src()}/>
        </li>
      )
    })
    return (
      <ul className="playbackSelector buttonRow controlSelector" onMouseOut={this.onMouseOut}>
        {buttons}
      </ul>
    )
  }
}

PlaybackSelector.propTypes = {
  mode: PropTypes.string.isRequired,
  onPlaybackChange: PropTypes.func.isRequired
}

export default PlaybackSelector