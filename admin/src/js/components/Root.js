import React, { PropTypes } from 'react'
import { AuthGlobals } from '../vendor_modules/redux-auth'

export default class Root extends React.Component {
  render () {

    const { connect, disconnect, dispatch, isAuthenticated, errorMessage } = this.props

    const children = React.Children.map(
      this.props.children,
      (child) => React.cloneElement(child, {
        connect,
        disconnect,
        dispatch,
        isAuthenticated,
        errorMessage
      })
    )

    return (
      <div id="root">
        { children }
      </div>      
    )
  }
}

Root.propTypes = {
  connect: PropTypes.func.isRequired,
  disconnect: PropTypes.func.isRequired,
  dispatch: PropTypes.func.isRequired,
  isAuthenticated: PropTypes.bool.isRequired,
  errorMessage: PropTypes.string
}
