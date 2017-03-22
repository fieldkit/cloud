
import React, { PropTypes } from 'react'
import { Link } from 'react-router'

class Navigation extends React.Component {
  render () {

    const {
      currentExpeditionID,
      currentPage
    } = this.props

    return (
      <ul className="navigation">
        <li 
          className={ 'navigation_link ' + (currentPage === 'map' ? 'active' : '') }
        >
          <Link
            to={ '/' + currentExpeditionID + '/map' }
          >
            Map
          </Link>
        </li>
        <li 
          className={ 'navigation_link ' + (currentPage === 'journal' ? 'active' : '') }
        >
          <Link
            to={ '/' + currentExpeditionID + '/journal' }
          >
            Journal
          </Link>
        </li>
      {/*
        <li 
          className={ 'navigation_link ' + (currentPage === 'about' ? 'active' : '') }
        >
          <Link
            to={ '/' + currentExpeditionID + '/about' }
          >
            About
          </Link>
        </li>
      */}
      </ul>
    )
  }

}

Navigation.propTypes = {

}

export default Navigation



