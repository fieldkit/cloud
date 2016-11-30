
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
          <ul className="sections">
            <li>
            <Link to={'/admin/' + id + '/uploader'}>Upload Data</Link>
            </li>
            <li>
            <Link to={'/admin/' + id + '/teams'}>Teams</Link>
            </li>
            <li>
            <Link to={'/admin/' + id + '/sources'}>Data Sources</Link>
            </li>
            <li>
            <Link to={'/admin/' + id + '/editor'}>Data Editor</Link>
            </li>
            <li>
            <Link to={'/admin/' + id + '/identity'}>Expedition Identity</Link>
            </li>
          </ul>
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
