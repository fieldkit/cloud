// @flow weak

import React, {Component} from 'react'

import './css/App.css'

import fieldkitLogoRed from './img/fieldkit-logo-red.svg'
import nationalGeographicLogoLong from './img/national-geographic-logo-long.png'
import conservifyLogo from './img/conservify_logo.png'
import ocrLogo from './img/ocr_logo.jpg'

export class App extends Component {
  render () {
    return (
      <div id="landing-page" className="page">

        <div className="slides">
          <div className="slide">
            <div className="content">
              <img src={nationalGeographicLogoLong} alt="" className="content_nat"/>
              <h1 className="content_title"><img className="content_title_img" src={fieldkitLogoRed} alt="Fieldkit logo"/></h1>
              <p className="content_sub">A one-click open platform for field researchers and explorers</p>

              <div id="mc_embed_signup">
              <form action="//nyc.us14.list-manage.com/subscribe/post?u=bedac4d8c1a0840ba87f528cd&amp;id=24d0f362c7" method="post" id="mc-embedded-subscribe-form" name="mc-embedded-subscribe-form" className="validate" target="_blank" noValidate>
                  <div id="mc_embed_signup_scroll">
                    <input type="email" name="EMAIL" className="email" id="mce-EMAIL" placeholder="email address" required/>
                    <div style={{ position: 'absolute', left: '-5000px' }} aria-hidden="true">
                      <input type="text" name="b_bedac4d8c1a0840ba87f528cd_24d0f362c7" tabIndex="-1" value=""/>
                    </div>
                    <div className="clear">
                      <input type="submit" value="Get notified when we launch" name="subscribe" id="mc-embedded-subscribe" className="button"/>
                    </div>
                  </div>
              </form>
              </div>
            </div>
          </div>
        </div>

        <footer>
          <ul className="footer_logos">
            <li className="footer_logo">
              <a href="http://conservify.org/" target="_blank" className="footer_logo_link">
                <img src={conservifyLogo} alt="Conservify" className="footer_logo_img"/>
              </a>
            </li>
            <li className="footer_logo">
              <a href="https://ocr.nyc/" target="_blank" className="footer_logo_link">
                <img src={ocrLogo} alt="The Office for Creative Research" className="footer_logo_img"/>
              </a>
            </li>
          </ul>
        </footer>
      </div>
    )
  }
}
