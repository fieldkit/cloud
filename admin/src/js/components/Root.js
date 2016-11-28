import React, { PropTypes } from 'react'

export default class Root extends React.Component {
  render () {

    const { connect } = this.props
    
    const children = React.Children.map(
      this.props.children,
      (child) => React.cloneElement(child, {
        connect: connect
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
  connect: PropTypes.func.isRequired
}
