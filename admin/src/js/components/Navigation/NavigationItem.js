
import React, { PropTypes } from 'react'
import { Link } from 'react-router'

class NavigationItem extends React.Component {

  render () {
    const { expedition, active, projectID } = this.props
    return (
      <li className={ active ? 'active' : ''}>
        <Link to={'/admin/' + projectID + '/' + expedition.get('id') + '/dashboard'}>
          <h4>{expedition.get('name') || 'New Expedition' }</h4>
        </Link>
        { 
          active &&
          <ul className="sections">
            <li>
            <Link to={'/admin/' + projectID + '/' + expedition.get('id') + '/uploader'}>Upload Data</Link>
            </li>
            <li>
            <Link to={'/admin/' + projectID + '/' + expedition.get('id') + '/teams'}>Teams</Link>
            </li>
            <li>
            <Link to={'/admin/' + projectID + '/' + expedition.get('id') + '/sources'}>Data Sources</Link>
            </li>
            <li>
            <Link to={'/admin/' + projectID + '/' + expedition.get('id') + '/editor'}>Data Editor</Link>
            </li>
            <li>
            <Link to={'/admin/' + projectID + '/' + expedition.get('id') + '/identity'}>Expedition Identity</Link>
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
  // active: PropTypes.bool
  // pathName: PropTypes.node.
}

export default NavigationItem
