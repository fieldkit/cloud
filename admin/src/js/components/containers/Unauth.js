// @flow weak

import React, { Component } from 'react'

import '../../../css/unauth.css'

export class Unauth extends Component {
  render () {
    return (
      <div className="unauth">
        { this.props.children }
      </div>
    )
  }
}
