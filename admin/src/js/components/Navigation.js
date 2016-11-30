import React, { PropTypes } from 'react'
import NavigationItem from './NavigationItem'

const Navigation = ({setPage}) => {
          // <NavigationItem setPage={setPage} active={pathName === '/data'}>Data</NavigationItem>
          // <NavigationItem active={pathName === '/share'}>Share</NavigationItem>
  return (
    <div id="header">
      <div id="navigation">
        <ul>
          <NavigationItem setPage={setPage} active={location.pathname.indexOf('/map') > -1}>Map</NavigationItem>
          <NavigationItem setPage={setPage} active={location.pathname.indexOf('/journal') > -1}>Journal</NavigationItem>
          <NavigationItem setPage={setPage} active={location.pathname.indexOf('/data') > -1}>Data</NavigationItem>
          <NavigationItem setPage={setPage} active={location.pathname.indexOf('/about') > -1}>About</NavigationItem>
        </ul>
      </div>
      <h1>INTO THE OKAVANGO</h1>
      <img id="logo" src="/static/img/logo.svg" alt="Into the Okavango"/>
    </div>
  )
}

Navigation.propTypes = {
  setPage: PropTypes.func.isRequired
}

export default Navigation