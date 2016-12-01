import React, {PropTypes} from 'react'
import { Link } from 'react-router'
import { EmailSignInForm } from "../vendor_modules/redux-auth/default-theme"
import autobind from 'autobind-decorator'
import { loginUser } from '../actions'


class SignInPage extends React.Component {
  constructor (props) {
    super(props)
    this.state = {

    }
  }

  @autobind
  onSignIn (event) {
    const email = this.refs.email
    const password = this.refs.password
    const creds = { email: email.value.trim(), password: password.value.trim() }
    this.props.dispatch(loginUser(creds))
    event.preventDefault()
    return false
  }

  render () {

    const { connect, errorMessage } = this.props

    return (
      <div id="signin-page" className="page">
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
          {/*
            <EmailSignInForm 
              endpoint={'http://localhost:3000/api/user/sign-in'}
              next={connect}
            />
          */}
          <form>
            <input type='text' ref='email' className="form-control" placeholder='Email'/>
            <input type='password' ref='password' className="form-control" placeholder='Password'/>
            <button onClick={this.onSignIn} className="btn btn-primary">
              Login
            </button>

            {errorMessage &&
              <p>{errorMessage}</p>
            }
          </form>

          <div onClick={connect}>
            <Link to={'/admin/okavango_16'}>(Fake sign in)</Link>
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
  connect: PropTypes.func,
  dispatch: PropTypes.func,
  isAuthenticated: PropTypes.bool.isRequired,
  errorMessage: PropTypes.string
}

export default SignInPage