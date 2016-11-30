import React, {PropTypes} from 'react'
import { Link } from 'react-router'
import { EmailSignInForm } from "redux-auth/default-theme";


class SignInPage extends React.Component {
  constructor (props) {
    super(props)
    this.state = {

    }
  }

  render () {

    const { connect } = this.props

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
          <EmailSignInForm 
            endpoint={'localhost:3000/signin'}
            next={connect}
          />
          <div onClick={connect}>
            <Link to={'/admin/okavango_16'}>Sign in</Link>
          </div>
          <p className="forgot-label">
            <Link to={'/forgot'}>Forgot your password?</Link>
          </p>
        </div>
      </div>
    )
  }

}

SignInPage.propTypes = {
  connect: PropTypes.func.isRequired
}

export default SignInPage