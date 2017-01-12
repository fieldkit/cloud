import React, {PropTypes} from 'react'
import { Link } from 'react-router'

import '../../scss/app.scss'

class LandingPage extends React.Component {
  constructor (props) {
    super(props)
    this.state = {

    }
  }

  render () {
    return (
      <div id="landing-page" className="page">
        <div 
          id="auth-panel"
          style={{
            position: 'absolute',
            right: 0
          }}
        >
          <Link to={'/signin'}>Sign in</Link>
          <Link to={'/signup'}>Sign up</Link>
        </div>
        <div className="content">
          <h1>FieldKit</h1>
          <Link to={'/signup'}>Get Started</Link>
        </div>
      </div>
    )
  }
}

LandingPage.propTypes = {

}

export default LandingPage