import React, {PropTypes} from 'react'
import { Link } from 'react-router'

import Signup from './Signup'

import '../../scss/app.scss'

class LandingPage extends React.Component {
  constructor (props) {
    super(props)
    this.state = {

    }
  }

  setSlidedState(value, ev) {
    ev.preventDefault()
    console.log("state")
    this.setState({
      slided : value
    })
  }

  render () {
    return (
      <div id="landing-page" className={"page landing-page " + (this.state.slided ? "slided" : "")} style={{'background-image':'url(/src/img/fieldkit-background.png)'}}>
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
              <a onClick={this.setSlidedState.bind(this, true)} href="#" className="navigation_link sign-up">Join us</a>
            </li>
          </ul>
        </nav>

        {(this.state.slided &&
          <a href="#"><i className="slide-back fa fa-long-arrow-up" onClick={this.setSlidedState.bind(this, false)} aria-hidden="true"></i></a>
        )}

        <div className="slide">
          <div className="content">
            <img src="/src/img/national-geographic-logo-long.png" alt="" className="content_nat"/>
            <h1 className="content_title"><img className="content_title_img" src="/src/img/fieldkit-logo-red.svg" alt="fieldkit"/></h1>
            <p className="content_sub">A one-click open platform for field researchers and explorers</p>
            <a onClick={this.setSlidedState.bind(this, true)} href="#" className="content_join">Join Us</a>
          </div>
        </div>

        <div className="slide">
          <div className="content">
            <Signup />
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