import React, {PropTypes} from 'react'
import { Link } from 'react-router'

import { EmailSignUpForm } from "redux-auth/default-theme";

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
          <EmailSignUpForm 
            endpoint={'localhost:3000/signup'}
          />
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

}

export default SignUpPage