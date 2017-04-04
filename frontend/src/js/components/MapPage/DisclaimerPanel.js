import React, { PropTypes } from 'react'
import { is } from 'immutable'

class DisclaimerPanel extends React.Component {

  render () {
    const {
    } = this.props
    
    return (
      <div className="disclaimer-panel">
        NOTE: Map images have been obtained from a third-party and do not reflect the editorial decisions of National Geographic.
      </div>
    )
  }

}

DisclaimerPanel.propTypes = {

}

export default DisclaimerPanel
