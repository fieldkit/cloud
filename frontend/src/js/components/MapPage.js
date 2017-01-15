
import React, { PropTypes } from 'react'
import ControlPanelContainer from '../containers/ControlPanelContainer'

class MapPage extends React.Component {

  constructor (props) {
    super(props)
    this.state = {
    }
  }

  render () {
    return (
      <div className="map-page page">
        <div className="page_content">
          <ControlPanelContainer/>
        </div>
      </div>
    )
  }

}

MapPage.propTypes = {

}

export default MapPage