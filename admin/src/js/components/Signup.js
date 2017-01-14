import React from 'react'

export default class Signup extends React.Component {
	render() {
		return (
			<div className="signup">
				<header className="signup_header">
					<h1 className="signup_title">Join us</h1>	
				</header>

				<form onSubmit={(ev) => ev.preventDefault()} className="signup_form">
					<div className="signup_content">
						<div className="signup_group">
							<label for="email" className="signup_label">E-mail</label>
							<input id="email" name="email" className="signup_input" type="email" placeholder="e-mail" />
						</div>

						<div className="signup_group">
							<label for="username" className="signup_label">Username</label>
							<input id="username" name="username" className="signup_input" type="text" placeholder="username" />
						</div>

						<div className="signup_group">
							<label for="password" className="signup_label">Password</label>
							<input id="password" name="password" className="signup_input" type="password" placeholder="password" />
						</div>

						<div className="signup_group">
							<label for="project" className="signup_label">Project's Name</label>
							<input id="project" name="project" className="signup_input" type="text"  placeholder="project's name"/>
						</div>
					</div>

					<footer class="signup_footer">
						<input className="signup_submit" type="submit" value="Submit"/>
					</footer>
				</form>
			</div>
		)
	}
}
