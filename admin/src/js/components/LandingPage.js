import React, {PropTypes} from 'react'
import { browserHistory, Link } from 'react-router'

import Signup from './Signup'
import Signin from './Signin'

import '../../scss/app.scss'

import fieldkitBackground from '../../img/fieldkit-background.png'
import fieldkitLogoRed from '../../img/fieldkit-logo-red.png'
import nationalGeographicLogoLong from '../../img/national-geographic-logo-long.png'
import conservifyLogo from '../../img/conservify_logo.png'
import ocrLogo from '../../img/ocr_logo.jpg'

class LandingPage extends React.Component {

  render () {

    const { 
      requestSignUp,
      requestSignIn
    } = this.props

    return (
      <div id="landing-page" className={"page landing-page " + (this.props.location.pathname === '/' ? '' : 'slided')} style={{'backgroundImage':'url(' + fieldkitBackground + ')'}}>

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
            <img src={nationalGeographicLogoLong} alt="" className="content_nat"/>
            <h1 className="content_title"><img className="content_title_img" src={fieldkitLogoRed} alt="fieldkit"/></h1>
            <p className="content_sub">A one-click open platform for field researchers and explorers</p>
            {/*
              <Link 
                to={'/signup'}
                className="content_join"
              >
                Join Us
              </Link>
            */}
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
                <img src={conservifyLogo} alt="" className="footer_logo_img"/>
              </a>
            </li>
            <li className="footer_logo">
              <a href="https://ocr.nyc/" target="_blank" className="footer_logo_link">
                <img src={ocrLogo} alt="" className="footer_logo_img"/>
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