
import React, { PropTypes } from 'react'
import Navigation from './Navigation'

import iconHamburger from '../../../img/icon-hamburger.svg'

class Header extends React.Component {
  render () {
    const {
      expeditionName,
      currentPage,
      currentExpeditionID,
      openExpeditionPanel
    } = this.props

    return (
      <div className="header">
        <h1 className="header_title">
          <img
            src={ '/' + iconHamburger }
            onClick={ openExpeditionPanel }
          />
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



