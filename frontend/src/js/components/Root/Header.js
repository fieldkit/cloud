
import React, { PropTypes } from 'react'
import Navigation from './Navigation'

class Header extends React.Component {

  constructor (props) {
    super(props)
    this.sate = {}
  }

  render () {

    const {
      expeditionName,
      currentPage,
      currentExpeditionID
    } = this.props

    return (
      <div className="header">
        <h1 className="header_title">
          { expeditionName }
        </h1>
        <Navigation
          currentExpeditionID={ currentExpeditionID }
          currentPage={ currentPage }
        />
      </div>
    )
  }

}

Header.propTypes = {

}

export default Header



