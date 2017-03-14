import React, { PropTypes } from 'react'

export default class Root extends React.Component {
  render () {
    return (
      <div id="root">
        { this.props.children }
      </div>      
    )
  }
}

Root.propTypes = {}
