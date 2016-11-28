import React, {PropTypes} from 'react'
import { Link } from 'react-router'

class SignInPage extends React.Component {
  constructor (props) {
    super(props)
    this.state = {

    }
  }

  render () {
    return (
      <div id="signup-page" className="page">
        <div id="header">
          <h1>
            <Link to={'/'}>FieldKit</Link>
          </h1>
        </div>
        <div 
          id="auth-panel"
          style={{
            position: 'absolute',
            right: 0
          }}
        >
          <Link to={'/signup'}>Sign up</Link>
        </div>
        <div className="content">
          <h1>Sign in</h1>
          <form>
            Login:<br/>
            Password:<br/>
            <Link to={'/admin/okavango_16'}>Sign in</Link>
          </form>
          <p className="forgot-label">
            <Link to={'/forgot'}>Forgot your password?</Link>
          </p>
        </div>
      </div>
    )
  }

}

SignInPage.propTypes = {

}

export default SignInPage