import React, { PropTypes } from 'react'
// import autobind from 'autobind-decorator'
import ReactPIXI, { TilingSprite, Sprite, Point } from 'react-pixi'

class PIXIStage extends React.Component {
  render () {
    var Stage = ReactPIXI.Stage
    return (
      <Stage width={this.props.width} height={this.props.height} transparent={true} backgroundColor={0x000000}>
        {this.props.children}
      </Stage>
    )
  }
}

PIXIStage.propTypes = {
  width: PropTypes.number.isRequired,
  height: PropTypes.number.isRequired,
}

export default PIXIStage
