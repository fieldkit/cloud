import React, {PropTypes} from 'react'
import { Link } from 'react-router'
import autobind from 'autobind-decorator'
// import { , loginUser, loginError } from '../actions'
import { browserHistory } from 'react-router'

class SignInPage extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
    }
  }

  @autobind
  onSubmit (event) {
    event.preventDefault()
    this.props.requestSignIn(this.refs.userName.value, this.refs.password.value)
    return false
  }

  render () {
    const { connect, errorMessage, fetching } = this.props
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
          <form>
            <input type='text' ref='userName' className="form-control" placeholder='Username'/>
            <input type='password' ref='password' className="form-control" placeholder='Password'/>
            <button onClick={this.onSubmit} className="btn btn-primary">
              Login
            </button>
            {fetching &&
              <span className="spinning-wheel-container"><div className="spinning-wheel"></div></span>
            }
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
  requestSignIn: PropTypes.func.isRequired,
  connect: PropTypes.func,
  errorMessage: PropTypes.string,
  fetching: PropTypes.bool.isRequired
}

export default SignInPage