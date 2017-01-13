import React, { PropTypes } from 'react'
import MapContainer from '../containers/MapContainer'

export default class Root extends React.Component {
  render () {
    return (
      <div id="root">
        <MapContainer/>
        { this.props.children }
      </div>      
    )
  }
}

Root.propTypes = {}
