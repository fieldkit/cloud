import React, {PropTypes} from 'react'
import { Link } from 'react-router'

import '../../scss/app.scss'

class LandingPage extends React.Component {
  constructor (props) {
    super(props)
    this.state = {

    }
  }

  render () {
    return (
      <div id="landing-page" className="page landing-page" style={{'background-image':'url(/src/img/fieldkit-background.jpg)'}}>
        <nav className="navigation">
          <ul className="navigation_links">
            <li className="navigation_item">
              <a href="" className="navigation_link">About</a>
            </li> 
            <li className="navigation_item">
              <a href="" className="navigation_link">Gallery</a>
            </li>
            <li className="navigation_item">
              <a href="" className="navigation_link">Log in</a>
            </li>
            <li className="navigation_item">
              <a href="" className="navigation_link sign-up">Join us</a>
            </li>
          </ul>
        </nav>

        <div className="content">
          <h1 className="content_title"><img className="content_title_img" src="/src/img/fieldkit-logo-red.svg" alt="fieldkit"/></h1>
          <p className="content_sub">A one-click open platform for field researchers and explorers</p>
          <a href="#" className="content_join">Join Us</a>
        </div>

        <footer className="footer">
          <ul className="footer_logos">
            <li className="footer_logo">
              <a href="http://conservify.org/" target="_blank" className="footer_logo_link">
                <img src="/src/img/conservify_logo.png" alt="" className="footer_logo_img"/>
              </a>
            </li>
            <li className="footer_logo">
              <a href="http://www.nationalgeographic.com/" target="_blank" className="footer_logo_link">
                <img src="/src/img/national-geographic-logo.png" alt="" className="footer_logo_img"/>
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