import React, {PropTypes} from 'react'
import { Link } from 'react-router'

class ForgotPasswordPage extends React.Component {
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
          <Link to={'/signup'}>Sign up</Link>
        </div>
        <div className="content">
          <h1>Forgot your password?</h1>
          <form>
            Forgot form
            <p className="forgot-label">
              We will send you a link to reset your password.
            </p>
            <div className="button">
              Reset password
            </div>
          </form>
        </div>
      </div>
    )
  }

}

ForgotPasswordPage.propTypes = {

}

export default ForgotPasswordPage