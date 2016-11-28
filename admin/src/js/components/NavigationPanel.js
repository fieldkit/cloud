import React, { PropTypes } from 'react'
import NavigationItem from './NavigationItem'
import { SignOutButton } from 'redux-auth/default-theme'
import { Link } from 'react-router'

class NavigationPanel extends React.Component {
          // <NavigationItem active={pathName === '/data'}>Data</NavigationItem>
          // <NavigationItem active={pathName === '/share'}>Share</NavigationItem>

  render () {

    const { expeditionID, expeditions, disconnect } = this.props

    const items = expeditions.map(expedition => {
      return <NavigationItem {...expedition} active={expeditionID === expedition.id} key={expedition.id} />
    })

    return (
      <div id="header">
        <h1>FieldKit</h1>
        <Link to={'/admin/profile'}>Adjany Costa</Link> 
        <SignOutButton
          endpoint={'localhost:3000/signout'}
          next={disconnect}
        >
          X
        </SignOutButton>
        <div id="navigation">
          <ul>
            {items}
          </ul>
        </div>
      </div>
    )
  }

}

NavigationPanel.propTypes = {
  disconnect: PropTypes.func.isRequired
}

export default NavigationPanel