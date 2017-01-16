
import React, { PropTypes } from 'react'
import ControlPanelContainer from '../containers/ControlPanelContainer'
import NotificationPanelContainer from '../containers/NotificationPanelContainer'

class MapPage extends React.Component {

  constructor (props) {
    super(props)
    this.state = {
    }
  }

  render () {
    return (
      <div className="map-page page">
        <ControlPanelContainer/>
        <NotificationPanelContainer/>
      </div>
    )
  }

}

MapPage.propTypes = {

}

export default MapPage