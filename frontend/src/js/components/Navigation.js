
import React, { PropTypes } from 'react'

class Navigation extends React.Component {

  constructor (props) {
    super(props)
    this.sate = {}
  }

  render () {

    return (
      <ul className="navigation">
        <li className="navigation_link active">
          Map
        </li>
        <li className="navigation_link">
          About
        </li>
      </ul>
    )
  }

}

Navigation.propTypes = {

}

export default Navigation



