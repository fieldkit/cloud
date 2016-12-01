import React, {PropTypes} from 'react'
import { Link } from 'react-router'
import autobind from 'autobind-decorator'
import { signupUser } from '../actions'

// import { EmailSignUpForm } from "../vendor_modules/redux-auth/default-theme";

class SignUpPage extends React.Component {
  constructor (props) {
    super(props)
    this.state = {

    }
  }

  @autobind
  onSignUp(event) {
    const email = this.refs.email
    const password = this.refs.password
    const creds = { email: email.value.trim(), password: password.value.trim() }
    this.props.dispatch(signupUser(creds))
    event.preventDefault()
    return false
  }

  render () {

    const { errorMessage } = this.props

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
          {/*
          <EmailSignUpForm 
            endpoint={'http://localhost:3000/api/user/sign-up'}
          />
          */}
          <form>
            <input type='text' ref='email' className="form-control" placeholder='Email'/>
            <input type='password' ref='password' className="form-control" placeholder='Password'/>
            <button onClick={this.onSignUp} className="btn btn-primary">
              Sign Up
            </button>
            {errorMessage &&
              <p>{errorMessage}</p>
            }
          </form>
          <p className="signin-label">
            Already have an account? <Link to={'/signin'}>Sign in</Link>
          </p>
        </div>
      </div>
    )
  }

}

/*

  endpoint: The key of the target provider service as represented in the endpoint configuration block.
  inputProps: An object containing the following attributes:
    email: An object that will override the email input component's default props.
    password: An object that will override the password input component's default props.
    passwordConfirmation: An object that will override the password confirmation input component's default props.
    submit: An object that will override the submit button component's default props.


*/

SignUpPage.propTypes = {
  connect: PropTypes.func,
  dispatch: PropTypes.func,
  isAuthenticated: PropTypes.bool.isRequired,
  errorMessage: PropTypes.string
}

export default SignUpPage