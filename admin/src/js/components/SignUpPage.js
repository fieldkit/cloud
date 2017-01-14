import React, {PropTypes} from 'react'
import { Link } from 'react-router'
import autobind from 'autobind-decorator'
import { signupUser, errorMessage } from '../actions'
import { browserHistory } from 'react-router'
import {FKApiClient} from '../api/api.js';

class SignUpPage extends React.Component {
  constructor (props) {
    super(props)
    this.state = {

    }
  }

  @autobind
  onSubmit (event) {
    event.preventDefault()
    this.props.requestSignUp(
      this.refs.email.value, 
      this.refs.userName.value,
      this.refs.password.value,
      this.refs.invite.value,
      this.refs.projectName.value
    )
    return false
  }

  render () {


    const { errorMessage, success, fetching } = this.props

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
            <input type='text' ref='invite' className="form-control" placeholder='Invite'/>
            <input type='text' ref='email' className="form-control" placeholder='Email'/>
            <input type='text' ref='userName' className="form-control" placeholder='User Name'/>
            <input type='password' ref='password' className="form-control" placeholder='Password'/>
            <input type='text' ref='projectName' className="form-control" placeholder='Project Name'/>
            <button onClick={this.onSubmit} className="btn btn-primary">
              Sign Up
            </button>
            {fetching &&
              <span className="spinning-wheel-container"><div className="spinning-wheel"></div></span>
            }
            {errorMessage &&
              <p>{errorMessage}</p>
            }
          </form>
          {!success &&
            <p className="signin-label">
              Already have an account? <Link to={'/signin'}>Sign in</Link>
            </p>
          }
          {success &&
            <p>
              You successfully signed up to FieldKit! <Link to={'/signin'}>Sign in</Link>
            </p>
          }
        </div>
      </div>
    )
  }
}

SignUpPage.propTypes = {
  requestSignUp: PropTypes.func.isRequired,
  errorMessage: PropTypes.string,
  success: PropTypes.bool,
  fetching: PropTypes.bool.isRequired
}

export default SignUpPage