import React, {PropTypes} from 'react'
import { Link } from 'react-router'

export default class Signin extends React.Component {

  constructor (props) {
    super(props)
    this.state = {}
    this.onSubmit = this.onSubmit.bind(this)
  } 

  onSubmit (event) {
    event.preventDefault()
    this.props.requestSignIn(
      this.refs.username.value, 
      this.refs.password.value
    )
    return false
  }

  render() {

    const {
      errors
    } = this.props

    if (errors) console.log('aga', errors)
    else console.log('empty')

    return (
      <div className="signin">
        <header className="signin_header">
          <h1 className="signin_title">Sign in</h1> 
        </header>

        <form onSubmit={this.onSubmit} className="signin_form">
          <div className="signin_content">
            <div className="signin_group">
              <label for="username" className="signin_label">Username</label>
              <input ref="username" id="username" name="username" className="signin_input" type="username" placeholder="" />
              {
                !!errors && !!errors.get && !!errors.get('username') &&
                errors.get('username').map((error, i) => {
                  return (
                    <p 
                      key={'errors-username-' + i}
                      className="errors"
                    >
                      {error}
                    </p>
                  )
                })
              }
            </div>
            <div className="signin_group">
              <label for="password" className="signin_label">Password</label>
              <input ref="password" id="password" name="password" className="signin_input" type="password" placeholder="" />
              {
                !!errors && !!errors.get && !!errors.get('password') &&
                errors.get('password').map((error, i) => {
                  return (
                    <p 
                      key={'errors-password-' + i}
                      className="errors"
                    >
                      {error}
                    </p>
                  )
                })
              }
            </div>
          </div>
          <footer class="signin_footer">
            {
              !!errors &&
              <p className="errors">
                Username or password invalid. Check your information and try again.
              </p>
            }
            <input className="signin_submit" type="submit" value="Submit"/>
          </footer>
        </form>
      </div>
    )
  }
}
