import React, { PropTypes } from 'react'
import { AuthGlobals } from 'redux-auth/default-theme'

export default class Root extends React.Component {
  render () {

    const { connect } = this.props
    
    const children = React.Children.map(
      this.props.children,
      (child) => React.cloneElement(child, {
        connect,
        disconnect
      })
    )

    return (
      <div id="root">
        <AuthGlobals />
        { children }
      </div>      
    )
  }
}

Root.propTypes = {
  connect: PropTypes.func.isRequired,
  disconnect: PropTypes.func.isRequired
}
