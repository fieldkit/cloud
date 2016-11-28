import React, { PropTypes } from 'react'
import NavigationItem from './NavigationItem'
import { Link } from 'react-router'

class NavigationPanel extends React.Component {
          // <NavigationItem active={pathName === '/data'}>Data</NavigationItem>
          // <NavigationItem active={pathName === '/share'}>Share</NavigationItem>

  render () {

    const { expeditionID, expeditions } = this.props

    const items = expeditions.map(expedition => {
      return <NavigationItem {...expedition} active={expeditionID === expedition.id} key={expedition.id} />
    })

    return (
      <div id="header">
        <h1>FieldKit</h1>
        <Link to={'/admin/profile'}>Adjany Costa</Link> <Link to={'/'}>X</Link>
        <div id="navigation">
          <ul>
            {items}
          </ul>
        </div>
      </div>
    )
  }

}

NavigationPanel.propTypes = {}

export default NavigationPanel