import React, { PropTypes } from 'react'
import TimelineContainer from '../containers/TimelineContainer'
import HeaderContainer from '../containers/HeaderContainer'
import MapContainer from '../containers/MapContainer'

export default class Root extends React.Component {
  render () {
    return (
      <div className="root">
        <MapContainer/>
        <div className="root_content">
          <HeaderContainer/>
          <TimelineContainer/>
          { this.props.children }
        </div>
      </div>
    )
  }
}

Root.propTypes = {}
