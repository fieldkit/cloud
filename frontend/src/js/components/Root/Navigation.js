
import React, { PropTypes } from 'react'
import { Link } from 'react-router'

class Navigation extends React.Component {
    render() {
        const { currentExpeditionID, currentPage } = this.props

        return (
            <ul className="navigation">
                <li className={ 'navigation_link ' + (currentPage === 'map' ? 'active' : '') }>
                    <Link to={ '/' + currentExpeditionID + '/map' }>
                        Map
                    </Link>
                </li>
            </ul>
        )
    }
}

Navigation.propTypes = {
}

export default Navigation
