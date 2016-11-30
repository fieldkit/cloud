import React, { PropTypes } from 'react'
import NavigationItem from './NavigationItem'
import { SignOutButton } from 'redux-auth/default-theme'
import { Link } from 'react-router'

class NavigationPanel extends React.Component {

  render () {
    const { expeditionID, expeditions, disconnect } = this.props

    const items = expeditions.map(expedition => {
      return <NavigationItem {...expedition} active={expeditionID === expedition.id} key={expedition.id} />
    })

    return (
      <div id="header">
        <div id="logo">
          <Link to="/admin">
            <img src="/src/img/fieldkit-logo.svg" alt="fieldkit logo" />
          </Link>
          <Link to={'/admin/profile'}>
            <img src="/src/img/profile-button.png" />
          </Link> 
          {/*
          <SignOutButton
            endpoint={'localhost:3000/signout'}
            next={disconnect}
          >
            X
          </SignOutButton>
          */}
        </div>
        <div id="navigation">
          <ul className="expeditions">
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