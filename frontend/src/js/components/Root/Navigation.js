
import React, { PropTypes } from 'react'
import { Link } from 'react-router'

class Navigation extends React.Component {
  render () {

    const {
      currentExpeditionID,
      currentProjectID,
      currentPage
    } = this.props
    
    let apiDomain = document.location.host.match(/fieldkit\.org/) ?
                        "https://api.fieldkit.org" :
                        "https://api.fieldkit.team"

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
        <li 
          className={ 'navigation_link ' + (currentPage === 'data' ? 'active' : '') }
        >
          <a href={ `${apiDomain}/projects/@/${currentProjectID}/expeditions/@/${currentExpeditionID}/documents` }>
            Data
          </a>
        </li>
      </ul>
    )
  }

}

Navigation.propTypes = {

}

export default Navigation



