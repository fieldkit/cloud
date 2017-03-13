import React, {PropTypes} from 'react'
import { Link } from 'react-router'

export default class Signup extends React.Component {

	constructor (props) {
		super(props)
		this.state = {}
		this.onSubmit = this.onSubmit.bind(this)
	} 

	onSubmit (event) {
    event.preventDefault()
    this.props.requestSignUp(
      this.refs.email.value, 
      this.refs.username.value,
      this.refs.password.value,
      this.refs.invite.value
    )
    return false
  }

	render() {

		const {
      errors
    } = this.props

		return (
			<div className="signup">
				<header className="signup_header">
					<h1 className="signup_title">Join us</h1>	
				</header>

				<form onSubmit={this.onSubmit} className="signup_form">
					<div className="signup_content">
						<div className="signup_group">
							<label for="email" className="signup_label">E-mail</label>
							<input ref="email" id="email" name="email" className="signup_input" type="email" placeholder="explorer@email.com" />
							{
	              !!errors && !!errors.get && !!errors.get('email') &&
	              errors.get('email').map((error, i) => {
	                return (
	                  <p 
	                    key={'errors-email-' + i}
	                    className="errors"
	                  >
	                    {error}
	                  </p>
	                )
	              })
	            }
						</div>

						<div className="signup_group">
							<label for="username" className="signup_label">Username</label>
							<input ref="username" id="username" name="username" className="signup_input" type="text" placeholder="explorer123" />
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

						<div className="signup_group">
							<label for="password" className="signup_label">Password</label>
							<input ref="password" id="password" name="password" className="signup_input" type="password" placeholder="correct horse battery staple" />
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

						<div className="signup_group">
							<label for="invite" className="signup_label">Invitation token</label>
							<input ref="invite" id="invite" name="invite" className="signup_input" type="text"  placeholder="E3NANDTJ3YXM5LNMGNZTF2373LAFTOCC"/>
							{
	              !!errors && !!errors.get && !!errors.get('invite') &&
	              errors.get('invite').map((error, i) => {
	                return (
	                  <p 
	                    key={'errors-invite-' + i}
	                    className="errors"
	                  >
	                    {error}
	                  </p>
	                )
	              })
	            }
						</div>
					</div>

					<footer class="signup_footer">
						{
		          !!errors &&
		          <p className="errors">
		            We found one or multiple errors. Please check your information above or try again later.
		          </p>
		        }
						<input className="signup_submit" type="submit" value="Submit"/>
					</footer>
				</form>
			</div>
		)
	}
}
