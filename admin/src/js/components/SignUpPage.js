import React, {PropTypes} from 'react'
import { Link } from 'react-router'

class SignUpPage extends React.Component {
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
          <Link to={'/signin'}>Sign in</Link>
        </div>
        <div className="content">
          <h1>Sign up</h1>
          <form>
            Sign up form
          </form>
          <p className="signin-label">
            Already have an account? <Link to={'/signin'}>Sign in</Link>
          </p>
        </div>
      </div>
    )
  }

}

SignUpPage.propTypes = {

}

export default SignUpPage