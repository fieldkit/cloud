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
    this.props.requestSignIn(this.refs.email.value, this.refs.password.value)
    return false
  }

  // async onSubmit (event) {
  //   console.log('login in as', this.refs.email.value, this.refs.password.value)
  //   event.preventDefault()
  //   try {
  //     await FKApiClient.get().login(this.refs.email, this.refs.password)
  //     // this.props.loginChanged()
  //     browserHistory.push('/admin/okavango_16')
  //   } catch (error) {

  //     console.log(error)

  //     if(error.response) {
  //       switch(error.response.status){
  //         case 429:
  //           this.props.dispatch(loginError('Try again later.'))
  //           break
  //         case 401:
  //           this.props.dispatch(loginError('Username or password incorrect.'))
  //           break
  //       }
  //     } else {
  //       this.props.dispatch(loginError('A server error occured.'))
  //     }
  //   }
  //   return false
  // }

  render () {

    const { connect, signinError } = this.props

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
            <button onClick={this.onSubmit} className="btn btn-primary">
              Login
            </button>

            {signinError &&
              <p>{signinError}</p>
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
  errorMessage: PropTypes.string
}

export default SignInPage