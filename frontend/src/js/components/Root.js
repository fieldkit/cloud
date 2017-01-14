import React, { PropTypes } from 'react'
import TimelineContainer from '../containers/TimelineContainer'
import MapContainer from '../containers/MapContainer'

export default class Root extends React.Component {
  render () {
    return (
      <div id="root">
        <MapContainer/>
        <div id="content">
          <TimelineContainer/>
          { this.props.children }
        </div>
      </div>
    )
  }
}

Root.propTypes = {}
