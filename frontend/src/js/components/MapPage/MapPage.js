
import React, { PropTypes } from 'react'
import ControlPanelContainer from '../../containers/common/ControlPanel'
import NotificationPanelContainer from '../../containers/MapPage/NotificationPanel'
import DisclaimerPanelContainer from '../../containers/MapPage/DisclaimerPanel'

class MapPage extends React.Component {
    render() {
        return (
            <div className="map-page page">
                <ControlPanelContainer/>
                <NotificationPanelContainer/>
                <DisclaimerPanelContainer/>
            </div>
        )
    }
}

MapPage.propTypes = {
}

export default MapPage
