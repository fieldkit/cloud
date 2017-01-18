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
            </div>
            <div className="signin_group">
              <label for="password" className="signin_label">Password</label>
              <input ref="password" id="password" name="password" className="signin_input" type="password" placeholder="" />
            </div>
          </div>
          <footer class="signin_footer">
            <input className="signin_submit" type="submit" value="Submit"/>
          </footer>
        </form>
      </div>
    )
  }
}
