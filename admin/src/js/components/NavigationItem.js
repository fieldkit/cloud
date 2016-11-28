
import React, { PropTypes } from 'react'
import { Link } from 'react-router'

class NavigationItem extends React.Component {

  render () {
    const { id, name, active } = this.props

    return (
      <li className={active ? 'active' : ''}>
        <Link to={'/admin/' + id + '/dashboard'}>
          <h4>{name}</h4>
        </Link>
        { 
          active &&
          <div>
            <Link to={'/admin/' + id + '/uploader'}>Upload Data</Link>
            <Link to={'/admin/' + id + '/teams'}>Teams</Link>
            <Link to={'/admin/' + id + '/sources'}>Data Sources</Link>
            <Link to={'/admin/' + id + '/editor'}>Data Editor</Link>
            <Link to={'/admin/' + id + '/identity'}>Expedition Identity</Link>
          </div>
        }
      </li>
    )
  }
}

NavigationItem.propTypes = {
  // children: PropTypes.node.isRequired,
  // setPage: PropTypes.func
  active: PropTypes.bool
  // pathName: PropTypes.node.
}

export default NavigationItem
