import React, { PropTypes } from 'react'

export default class Root extends React.Component {
  render () {

    // const { connect, disconnect, dispatch, isAuthenticated, signupError, signinError } = this.props

    // const children = React.Children.map(
    //   this.props.children,
    //   (child) => React.cloneElement(child, {
    //     connect,
    //     disconnect,
    //     dispatch,
    //     isAuthenticated,
    //     signupError,
    //     signinError,
    //   })
    // )

    return (
      <div id="root">
        { this.props.children }
      </div>      
    )
  }
}

Root.propTypes = {}
