
// import React, { PropTypes } from 'react'
import React from 'react'
import NotificationPanelContainer from '../containers/NotificationPanelContainer'
import ControlPanelContainer from '../containers/ControlPanelContainer.js'

class MapPage extends React.Component {

  render () {
    var height = {height: window.innerWidth > 768 ? window.innerHeight - 100 : window.innerHeight - 120}
    return (
      <div className='page' id="mapPage" style={height}>
        <ControlPanelContainer/>
        <NotificationPanelContainer/>
      </div>
    )
  }
}

// MapPage.propTypes = {
//   active : PropTypes.bool.isRequired
// }
export default MapPage
