import React, {PropTypes} from 'react'
import { browserHistory, Link } from 'react-router'

import Signup from './Signup'
import Signin from './Signin'

import '../../scss/app.scss'

class LandingPage extends React.Component {

  render () {

    const { 
      requestSignUp,
      requestSignIn
    } = this.props

    return (
      <div id="landing-page" className={"page landing-page " + (this.props.location.pathname === '/' ? '' : 'slided')} style={{'backgroundImage':'url(/src/img/fieldkit-background.png)'}}>
        <nav className="navigation">
          <ul className="navigation_links">
            <li className="navigation_item">
              <a href="" className="navigation_link">About</a>
            </li> 
            <li className="navigation_item">
              <a href="" className="navigation_link">Gallery</a>
            </li>
            <li className="navigation_item">
              <Link
                to={'/signin'}
                className="navigation_link"
              >
                Log in
              </Link>
            </li>
            <li className="navigation_item">
              <Link
                to={'/signup'}
                className="navigation_link sign-up"
              >
                Join us
              </Link>
            </li>
          </ul>
        </nav>

        {(this.props.location.pathname !== '/' &&
          <Link to={'/'}>
            <i 
              className="slide-back fa fa-long-arrow-up"
              aria-hidden="true"
            ></i>
          </Link>
        )}

        <div className="slide">
          <div className="content">
            <img src="/src/img/national-geographic-logo-long.png" alt="" className="content_nat"/>
            <h1 className="content_title"><img className="content_title_img" src="/src/img/fieldkit-logo-red.svg" alt="fieldkit"/></h1>
            <p className="content_sub">A one-click open platform for field researchers and explorers</p>
            <Link 
              to={'/signup'}
              className="content_join"
            >
              Join Us
            </Link>
          </div>
        </div>

        <div className="slide">
          <div className="content">
            {
              this.props.location.pathname === '/signup' && 
              <Signup requestSignUp={requestSignUp} />
            }
            {
              this.props.location.pathname === '/signin' && 
              <Signin requestSignIn={requestSignIn} />
            }
          </div>
        </div>

        <footer className="footer">
          <ul className="footer_logos">
            <li className="footer_logo">
              <a href="http://conservify.org/" target="_blank" className="footer_logo_link">
                <img src="/src/img/conservify_logo.png" alt="" className="footer_logo_img"/>
              </a>
            </li>
            <li className="footer_logo">
              <a href="https://ocr.nyc/" target="_blank" className="footer_logo_link">
                <img src="/src/img/ocr_logo.jpg" alt="" className="footer_logo_img"/>
              </a>
            </li>
          </ul>
        </footer>
      </div>
    )
  }
}

LandingPage.propTypes = {

}

export default LandingPage