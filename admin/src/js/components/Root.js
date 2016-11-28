import React, { PropTypes } from 'react'

export default class Root extends React.Component {
  render () {
    
    const { children } = this.props

    return (
      <div id="root">
        { children }
      </div>      
    )
  }
}

Root.propTypes = {
}
