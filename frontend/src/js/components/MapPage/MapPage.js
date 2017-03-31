
import React, { PropTypes } from 'react'
import ControlPanelContainer from '../../containers/common/ControlPanel'
import NotificationPanelContainer from '../../containers/MapPage/NotificationPanel'

class MapPage extends React.Component {
  render () {
    const {
      documents
    } = this.props

    return (
      <div className="map-page page">
        {
          !!documents && documents.size > 0 &&
          <ControlPanelContainer/>
        }
        <NotificationPanelContainer/>
      </div>
    )
  }

}

MapPage.propTypes = {

}

export default MapPage